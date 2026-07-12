package tests

import (
	"context"
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var testUserSequence int64

type userIntegrationSubject struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	timeout time.Duration
}

type clerkTestUserData struct {
	Email        string
	UpdatedEmail string
	Username     string
	UpdatedName  string
}

func newUserIntegrationSubject(t *testing.T) *userIntegrationSubject {
	t.Helper()

	require.NoError(t, testdb.LoadEnv())
	if !integrationTestsEnabled() {
		t.Skip("set CLERK_INTEGRATION_TESTS=1 to run real Clerk + DB integration tests")
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if strings.TrimSpace(clerkSecretKey) == "" {
		t.Fatal("CLERK_SECRET_KEY must be set for real Clerk integration tests")
	}
	clerk.SetKey(clerkSecretKey)

	pool := testdb.Open(t)

	repos := repository.NewRepository(pool)
	return &userIntegrationSubject{
		pool:    pool,
		queries: repos.Queries,
		timeout: webhookWaitTimeout(),
	}
}

func integrationTestsEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("CLERK_INTEGRATION_TESTS")))
	return value == "1" || value == "true" || value == "yes"
}

func envOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func webhookWaitTimeout() time.Duration {
	value := strings.TrimSpace(os.Getenv("CLERK_WEBHOOK_WAIT_TIMEOUT"))
	if value == "" {
		return 2 * time.Minute
	}

	timeout, err := time.ParseDuration(value)
	if err == nil {
		return timeout
	}

	seconds, err := strconv.Atoi(value)
	if err == nil {
		return time.Duration(seconds) * time.Second
	}

	return 45 * time.Second
}

func uniqueClerkTestUser(t *testing.T) clerkTestUserData {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&testUserSequence, 1), 10)
	emailPrefix := "integration.user"

	return clerkTestUserData{
		Email:        emailPrefix + "+clerk_test_" + suffix + "@example.com",
		UpdatedEmail: emailPrefix + "+clerk_updated_" + suffix + "@example.com",
		Username:     "integration_user_" + suffix,
		UpdatedName:  "integration_updated_" + suffix,
	}
}

func createClerkUser(t *testing.T, ctx context.Context, userData clerkTestUserData) *clerk.User {
	t.Helper()

	requireLocalWebhookEndpointReachable(t)

	createdUser, err := user.Create(ctx, &user.CreateParams{
		EmailAddresses:          &[]string{userData.Email},
		Username:                clerk.String(userData.Username),
		SkipPasswordRequirement: clerk.Bool(true),
		SkipLegalChecks:         clerk.Bool(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, createdUser.ID)

	return createdUser
}

func requireLocalWebhookEndpointReachable(t *testing.T) {
	t.Helper()

	baseURL := strings.TrimRight(envOrDefault("CLERK_LOCAL_BACKEND_URL", "http://localhost:"+envOrDefault("PORT", "8000")), "/")
	parsedURL, err := url.Parse(baseURL)
	require.NoError(t, err)

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		switch parsedURL.Scheme {
		case "https":
			host += ":443"
		default:
			host += ":80"
		}
	}

	connection, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		t.Fatalf(
			"local backend is not reachable at %q; start the API with `go run ./cmd` before running real Clerk integration tests",
			baseURL,
		)
	}
	require.NoError(t, connection.Close())
}

func cleanupLocalUser(t *testing.T, queries *db.Queries, userID string) {
	t.Helper()

	err := queries.DeleteUserByClerkID(context.Background(), userID)
	require.NoError(t, err)
}

func countLocalUsersByEmail(t *testing.T, subject *userIntegrationSubject, email string) int64 {
	t.Helper()

	var count int64
	err := subject.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	require.NoError(t, err)

	return count
}

func requireUniqueViolation(t *testing.T, err error) {
	t.Helper()

	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.ErrorAs(t, err, &pgErr)
	require.Equal(t, "23505", pgErr.Code)
}

func waitForUser(t *testing.T, subject *userIntegrationSubject, userID string) db.User {
	t.Helper()

	deadline := time.Now().Add(subject.timeout)
	for time.Now().Before(deadline) {
		user, err := subject.queries.GetUserByClerkID(context.Background(), userID)
		if err == nil {
			return user
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			require.NoError(t, err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf(
		"timed out after %s waiting for Clerk webhook to create user %q in the local database; check Clerk webhook delivery attempts, the public forwarded URL, webhook signing secret, and that backend/test use the same PostgreSQL database",
		subject.timeout,
		userID,
	)
	return db.User{}
}

func waitForUserEmail(t *testing.T, subject *userIntegrationSubject, userID string, email string) db.User {
	t.Helper()

	deadline := time.Now().Add(subject.timeout)
	for time.Now().Before(deadline) {
		user, err := subject.queries.GetUserByClerkID(context.Background(), userID)
		require.NoError(t, err)
		if user.Email == email {
			return user
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf(
		"timed out after %s waiting for Clerk webhook to update user %q email to %q in the local database",
		subject.timeout,
		userID,
		email,
	)
	return db.User{}
}

func waitForUserUsername(t *testing.T, subject *userIntegrationSubject, userID string, username string) db.User {
	t.Helper()

	deadline := time.Now().Add(subject.timeout)
	for time.Now().Before(deadline) {
		user, err := subject.queries.GetUserByClerkID(context.Background(), userID)
		require.NoError(t, err)
		if user.Username == username {
			return user
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf(
		"timed out after %s waiting for Clerk webhook to update user %q username to %q in the local database",
		subject.timeout,
		userID,
		username,
	)
	return db.User{}
}

func waitForUserDeleted(t *testing.T, subject *userIntegrationSubject, userID string) {
	t.Helper()

	deadline := time.Now().Add(subject.timeout)
	for time.Now().Before(deadline) {
		_, err := subject.queries.GetUserByClerkID(context.Background(), userID)
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		require.NoError(t, err)

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf(
		"timed out after %s waiting for Clerk webhook to delete user %q from the local database",
		subject.timeout,
		userID,
	)
}
