package tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var roomIntegrationSequence int64

type roomIntegrationSubject struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

type roomTestUser struct {
	ID       string
	Username string
	Email    string
}

func newRoomIntegrationSubject(t *testing.T) *roomIntegrationSubject {
	t.Helper()
	loadRoomTestEnv(t)

	pool, err := repository.NewPostgresDB(context.Background(), roomIntegrationPostgresConfig())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	repos := repository.NewRepository(pool)
	return &roomIntegrationSubject{pool: pool, queries: repos.Queries}
}

func loadRoomTestEnv(t *testing.T) {
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

func roomIntegrationPostgresConfig() config.PostgresDBConfig {
	return config.PostgresDBConfig{
		Host:     roomEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     roomEnvOrDefault("POSTGRES_PORT", "5432"),
		Username: roomEnvOrDefault("POSTGRES_USER", "miskatonic_user"),
		Password: roomEnvOrDefault("POSTGRES_PASSWORD", "miskatonic_password"),
		DBName:   roomEnvOrDefault("POSTGRES_DB", "miskatonic_lab"),
		SSLMode:  roomEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func roomEnvOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func createRoomTestUser(t *testing.T, subject *roomIntegrationSubject) roomTestUser {
	t.Helper()

	suffix := uniqueRoomIntegrationSuffix()
	user := roomTestUser{
		ID:       "room_integration_user_" + suffix,
		Username: "room_integration_" + suffix,
		Email:    "room.integration+" + suffix + "@example.com",
	}

	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupRoomTestUser(t, subject.queries, user.ID)
	})

	return user
}

func cleanupRoomTestUser(t *testing.T, queries *db.Queries, userID string) {
	t.Helper()
	_ = queries.DeleteUserByClerkID(context.Background(), userID)
}

func createRoomTestRoom(t *testing.T, subject *roomIntegrationSubject, ownerID string) db.Room {
	t.Helper()

	room, err := subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      ownerID,
		MaxPlayers:   4,
		InviteToken:  "invite_" + uniqueRoomIntegrationSuffix(),
		PasswordHash: "test_password_hash",
	})
	require.NoError(t, err)

	return room
}

func addRoomTestMember(t *testing.T, subject *roomIntegrationSubject, roomID pgtype.UUID, userID string, role string) db.RoomMember {
	t.Helper()

	member, err := subject.queries.AddMember(context.Background(), db.AddMemberParams{
		RoomID: roomID,
		UserID: userID,
		Role:   role,
	})
	require.NoError(t, err)

	return member
}

func createRoomTestCharacter(t *testing.T, subject *roomIntegrationSubject, userID string) db.Character {
	t.Helper()

	character, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{
		UserID: userID,
		Name:   "Room Test Investigator",
	})
	require.NoError(t, err)

	return character
}

func setRoomLastActivityAt(t *testing.T, subject *roomIntegrationSubject, roomID pgtype.UUID, activityAt time.Time) {
	t.Helper()

	_, err := subject.pool.Exec(
		context.Background(),
		"UPDATE rooms SET last_activity_at = $1, updated_at = $1 WHERE id = $2",
		activityAt,
		roomID,
	)
	require.NoError(t, err)
}

func requireRoomLastActivityAfter(
	t *testing.T,
	subject *roomIntegrationSubject,
	roomID pgtype.UUID,
	userID string,
	activityAt time.Time,
) {
	t.Helper()

	room, err := subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{
		ID:     roomID,
		UserID: userID,
	})
	require.NoError(t, err)
	require.Truef(
		t,
		room.LastActivityAt.Time.After(activityAt),
		"expected room activity %s to be after %s",
		room.LastActivityAt.Time,
		activityAt,
	)
}

func uniqueRoomIntegrationSuffix() string {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&roomIntegrationSequence, 1), 10)
	}

	return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&roomIntegrationSequence, 1), 10) + "_" + hex.EncodeToString(randomBytes)
}

func requireRoomPostgresErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.ErrorAs(t, err, &pgErr)
	require.Equal(t, code, pgErr.Code)
}

func roomTestUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}
	return uuid
}
