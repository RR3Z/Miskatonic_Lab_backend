package tests

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func newE2ESubject(t *testing.T) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)
	require.NotNil(t, suiteE2EClerkFixture, "E2E Clerk fixture is not initialized")
	return newE2ESubjectForIdentity(t, suiteE2EClerkFixture.primary)
}

func newSecondE2ESubject(t *testing.T) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)
	require.NotNil(t, suiteE2EClerkFixture, "E2E Clerk fixture is not initialized")
	return newE2ESubjectForIdentity(t, suiteE2EClerkFixture.secondary)
}

func newE2ESubjectForIdentity(t *testing.T, identity e2eClerkIdentity) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)

	pool := testdb.Open(t)

	queries := repository.NewRepository(pool).Queries
	cleanupUser := ensureLocalE2EUser(t, queries, identity.userID)
	t.Cleanup(cleanupUser)

	return &e2eSubject{
		baseURL:  e2eBaseURL(),
		identity: identity,
		userID:   identity.userID,
		client:   &http.Client{Timeout: 10 * time.Second},
		pool:     pool,
		queries:  queries,
	}
}

func requireE2EEnabled(t *testing.T) {
	t.Helper()
	if !e2eEnabled() {
		t.Skip("set E2E_TESTS=1 to run live backend E2E tests")
	}
}

func e2eEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("E2E_TESTS")))
	return value == "1" || value == "true" || value == "yes"
}

func loadE2EEnv(t *testing.T) {
	t.Helper()

	require.NoError(t, testdb.LoadEnv())
}

func e2eBaseURL() string {
	value := strings.TrimSpace(os.Getenv("E2E_BASE_URL"))
	if value == "" {
		value = "http://localhost:" + e2eEnvOrDefault("PORT", "8000")
	}
	return strings.TrimRight(value, "/")
}

func ensureLocalE2EUser(t *testing.T, queries *db.Queries, userID string) func() {
	t.Helper()

	_, err := queries.GetUserByClerkID(context.Background(), userID)
	if err == nil {
		return func() {}
	}
	require.True(t, errors.Is(err, pgx.ErrNoRows), "unexpected user lookup error: %v", err)

	hash := e2eHash(userID)
	_, err = queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:       userID,
		Username: "e2e_" + hash,
		Email:    "e2e+" + hash + "@example.com",
	})
	require.NoError(t, err)

	return func() {
		_ = queries.DeleteUserByClerkID(context.Background(), userID)
	}
}

func e2eEnvOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
