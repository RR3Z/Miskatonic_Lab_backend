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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var diceIntegrationSequence int64

type diceIntegrationSubject struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

type diceTestUser struct {
	ID       string
	Username string
	Email    string
}

func newDiceIntegrationSubject(t *testing.T) *diceIntegrationSubject {
	t.Helper()
	loadDiceTestEnv(t)

	pool, err := repository.NewPostgresDB(context.Background(), diceIntegrationPostgresConfig())
	require.NoError(t, err)

	repos := repository.NewRepository(pool)
	t.Cleanup(func() {
		pool.Close()
	})

	return &diceIntegrationSubject{
		pool:    pool,
		queries: repos.Queries,
	}
}

func loadDiceTestEnv(t *testing.T) {
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

func diceIntegrationPostgresConfig() config.PostgresDBConfig {
	return config.PostgresDBConfig{
		Host:     diceEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     diceEnvOrDefault("POSTGRES_PORT", "5432"),
		Username: diceEnvOrDefault("POSTGRES_USER", "miskatonic_user"),
		Password: diceEnvOrDefault("POSTGRES_PASSWORD", "miskatonic_password"),
		DBName:   diceEnvOrDefault("POSTGRES_DB", "miskatonic_lab"),
		SSLMode:  diceEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func diceEnvOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func createDiceTestUser(t *testing.T, subject *diceIntegrationSubject) diceTestUser {
	t.Helper()

	suffix := uniqueDiceIntegrationSuffix()
	testUser := diceTestUser{
		ID:       "dice_integration_user_" + suffix,
		Username: "dice_integration_" + suffix,
		Email:    "dice.integration+" + suffix + "@example.com",
	}

	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Email:     testUser.Email,
		AvatarUrl: nil,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupDiceTestUser(t, subject.queries, testUser.ID)
	})

	return testUser
}

func createDiceTestCharacter(t *testing.T, subject *diceIntegrationSubject, userID string) db.Character {
	t.Helper()

	playerName := "Roger"
	occupation := "Antiquarian"
	age := int16(37)
	sex := "male"
	residence := "Arkham"
	birthplace := "Boston"

	character, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{
		UserID:     userID,
		Name:       "Dr. Henry Armitage",
		PlayerName: &playerName,
		Occupation: &occupation,
		Age:        &age,
		Sex:        &sex,
		Residence:  &residence,
		Birthplace: &birthplace,
	})
	require.NoError(t, err)

	return character
}

func cleanupDiceTestUser(t *testing.T, queries *db.Queries, userID string) {
	t.Helper()

	err := queries.DeleteUserByClerkID(context.Background(), userID)
	require.NoError(t, err)
}

func uniqueDiceIntegrationSuffix() string {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&diceIntegrationSequence, 1), 10)
	}

	return strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + strconv.FormatInt(atomic.AddInt64(&diceIntegrationSequence, 1), 10) + "_" + hex.EncodeToString(randomBytes)
}

func createTestDiceRollParams(userID string, characterID pgtype.UUID) db.CreateDiceRollParams {
	return db.CreateDiceRollParams{
		UserID:      userID,
		CharacterID: characterID,
		Expression:  "2d6+3",
		Result:      10,
		Details:     []byte(`[{"type":"dice","sides":6,"rolls":[3,4]},{"type":"modifier","value":3}]`),
	}
}

func diceTestUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}
	return uuid
}
