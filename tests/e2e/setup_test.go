package tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func newE2ESubject(t *testing.T) *e2eSubject {
	t.Helper()
	requireE2EEnabled(t)
	loadE2EEnv(t)

	token := normalizedE2EToken(t)
	userID := subjectFromJWT(t, token)

	pool, err := repository.NewPostgresDB(context.Background(), e2ePostgresConfig())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

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

	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			require.NoError(t, godotenv.Load(envPath))
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func e2eBaseURL() string {
	value := strings.TrimSpace(os.Getenv("E2E_BASE_URL"))
	if value == "" {
		value = "http://localhost:" + e2eEnvOrDefault("PORT", "8000")
	}
	return strings.TrimRight(value, "/")
}

func normalizedE2EToken(t *testing.T) string {
	t.Helper()
	token := strings.TrimSpace(os.Getenv("E2E_AUTH_TOKEN"))
	require.NotEmpty(t, token, "E2E_AUTH_TOKEN must be set when E2E_TESTS=1")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")
	require.NotEmpty(t, token, "E2E_AUTH_TOKEN must contain a token")
	return token
}

func subjectFromJWT(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3, "E2E_AUTH_TOKEN must be a JWT")

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)

	var claims struct {
		Subject string `json:"sub"`
	}
	require.NoError(t, json.Unmarshal(payload, &claims))
	require.NotEmpty(t, strings.TrimSpace(claims.Subject), "E2E_AUTH_TOKEN JWT must contain sub")
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

func e2ePostgresConfig() config.PostgresDBConfig {
	return config.PostgresDBConfig{
		Host:     e2eEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     e2eEnvOrDefault("POSTGRES_PORT", "5432"),
		Username: e2eEnvOrDefault("POSTGRES_USER", "miskatonic_user"),
		Password: e2eEnvOrDefault("POSTGRES_PASSWORD", "miskatonic_password"),
		DBName:   e2eEnvOrDefault("POSTGRES_DB", "miskatonic_lab"),
		SSLMode:  e2eEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func e2eEnvOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
