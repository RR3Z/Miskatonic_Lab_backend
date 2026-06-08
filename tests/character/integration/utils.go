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

var characterIntegrationSequence int64

type characterIntegrationSubject struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

type characterTestUser struct {
	ID       string
	Username string
	Email    string
}

func newCharacterIntegrationSubject(t *testing.T) *characterIntegrationSubject {
	t.Helper()
	loadCharacterTestEnv(t)

	pool, err := repository.NewPostgresDB(context.Background(), characterIntegrationPostgresConfig())
	require.NoError(t, err)

	repos := repository.NewRepository(pool)
	t.Cleanup(func() {
		pool.Close()
	})

	return &characterIntegrationSubject{
		pool:    pool,
		queries: repos.Queries,
	}
}

func loadCharacterTestEnv(t *testing.T) {
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

func characterIntegrationPostgresConfig() config.PostgresDBConfig {
	return config.PostgresDBConfig{
		Host:     characterEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     characterEnvOrDefault("POSTGRES_PORT", "5432"),
		Username: characterEnvOrDefault("POSTGRES_USER", "miskatonic_user"),
		Password: characterEnvOrDefault("POSTGRES_PASSWORD", "miskatonic_password"),
		DBName:   characterEnvOrDefault("POSTGRES_DB", "miskatonic_lab"),
		SSLMode:  characterEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func characterEnvOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	return value
}

func createCharacterTestUser(t *testing.T, subject *characterIntegrationSubject) characterTestUser {
	t.Helper()

	suffix := uniqueCharacterIntegrationSuffix()
	testUser := characterTestUser{
		ID:       "character_integration_user_" + suffix,
		Username: "character_integration_" + suffix,
		Email:    "character.integration+" + suffix + "@example.com",
	}

	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Email:     testUser.Email,
		AvatarUrl: nil,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupCharacterTestUser(t, subject.queries, testUser.ID)
	})

	return testUser
}

func createCharacterTestCharacter(t *testing.T, subject *characterIntegrationSubject, userID string) db.Character {
	t.Helper()

	character, err := subject.queries.CreateCharacter(context.Background(), testCreateCharacterParams(userID))
	require.NoError(t, err)

	return character
}

func cleanupCharacterTestUser(t *testing.T, queries *db.Queries, userID string) {
	t.Helper()

	err := queries.DeleteUserByClerkID(context.Background(), userID)
	require.NoError(t, err)
}

func uniqueCharacterIntegrationSuffix() string {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&characterIntegrationSequence, 1), 10)
	}

	return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&characterIntegrationSequence, 1), 10) + "_" + hex.EncodeToString(randomBytes)
}

func testCreateCharacterParams(userID string) db.CreateCharacterParams {
	playerName := "Roger"
	occupation := "Antiquarian"
	age := int16(37)
	sex := "male"
	residence := "Arkham"
	birthplace := "Boston"

	return db.CreateCharacterParams{
		UserID:     userID,
		Name:       "Dr. Henry Armitage",
		PlayerName: &playerName,
		Occupation: &occupation,
		Age:        &age,
		Sex:        &sex,
		Residence:  &residence,
		Birthplace: &birthplace,
	}
}

func requirePostgresErrorCode(t *testing.T, err error, code string) {
	t.Helper()

	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.ErrorAs(t, err, &pgErr)
	require.Equal(t, code, pgErr.Code)
}

func characterTestUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}

	return uuid
}
