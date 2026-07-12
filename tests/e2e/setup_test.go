package tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	return newE2ESubjectFromEnv(t, "E2E_AUTH_TOKEN")
}

func newSecondE2ESubject(t *testing.T) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)
	if strings.TrimSpace(os.Getenv("E2E_SECOND_AUTH_TOKEN")) == "" {
		t.Skip("set E2E_SECOND_AUTH_TOKEN to run real multi-user E2E tests")
	}
	return newE2ESubjectFromEnv(t, "E2E_SECOND_AUTH_TOKEN")
}

func newE2ESubjectFromEnv(t *testing.T, tokenEnv string) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)

	token := normalizedE2EToken(t, tokenEnv)
	userID := subjectFromJWT(t, token)

	pool := testdb.Open(t)

	queries := repository.NewRepository(pool).Queries
	cleanupUser := ensureLocalE2EUser(t, queries, userID)
	t.Cleanup(cleanupUser)

	return &e2eSubject{
		baseURL: e2eBaseURL(),
		token:   token,
		userID:  userID,
		client:  &http.Client{Timeout: 10 * time.Second},
		pool:    pool,
		queries: queries,
	}
}

func requireE2EEnabled(t *testing.T) {
	t.Helper()
	value := strings.ToLower(strings.TrimSpace(os.Getenv("E2E_TESTS")))
	if value != "1" && value != "true" && value != "yes" {
		t.Skip("set E2E_TESTS=1 to run live backend E2E tests")
	}
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

func normalizedE2EToken(t *testing.T, tokenEnv string) string {
	t.Helper()
	token := strings.TrimSpace(os.Getenv(tokenEnv))
	require.NotEmpty(t, token, "%s must be set when E2E_TESTS=1", tokenEnv)
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")
	require.NotEmpty(t, token, "%s must contain a token", tokenEnv)
	return token
}

func subjectFromJWT(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3, "E2E auth token must be a JWT")

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)

	var claims struct {
		Subject string `json:"sub"`
	}
	require.NoError(t, json.Unmarshal(payload, &claims))
	require.NotEmpty(t, strings.TrimSpace(claims.Subject), "E2E auth token JWT must contain sub")
	return claims.Subject
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
