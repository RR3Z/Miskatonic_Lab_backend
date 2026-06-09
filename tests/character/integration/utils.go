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
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
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

type characterEventRecorder struct {
	events []events.Event
}

func (r *characterEventRecorder) Publish(_ context.Context, event events.Event) {
	r.events = append(r.events, event)
}

func requireLastCharacterEvent[T events.Event](t *testing.T, recorder *characterEventRecorder) T {
	t.Helper()

	require.NotEmpty(t, recorder.events)
	event, ok := recorder.events[len(recorder.events)-1].(T)
	require.True(t, ok)

	return event
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

func characterInt16(value int16) *int16 {
	return &value
}

func characterString(value string) *string {
	return &value
}

func requireCharacteristicValue(t *testing.T, actual *int16, expected int16) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func requireDerivedStatValue(t *testing.T, actual *int16, expected int16) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func requireDerivedStatString(t *testing.T, actual *string, expected string) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func createSkillTestCategory(t *testing.T, subject *characterIntegrationSubject, name string) (pgtype.UUID, string) {
	t.Helper()

	var id pgtype.UUID
	uniqueName := name + " " + uniqueCharacterIntegrationSuffix()
	err := subject.pool.QueryRow(context.Background(),
		"INSERT INTO skills_categories (name) VALUES ($1) RETURNING id",
		uniqueName,
	).Scan(&id)
	require.NoError(t, err)

	return id, uniqueName
}

func createSkillTestSpecialty(t *testing.T, subject *characterIntegrationSubject, name string, description string, baseValue int16) (pgtype.UUID, string) {
	t.Helper()

	var id pgtype.UUID
	uniqueName := name + " " + uniqueCharacterIntegrationSuffix()
	err := subject.pool.QueryRow(context.Background(),
		"INSERT INTO skills_specialties (name, description, base_value) VALUES ($1, $2, $3) RETURNING id",
		uniqueName,
		description,
		baseValue,
	).Scan(&id)
	require.NoError(t, err)

	return id, uniqueName
}

func testCreateSkillParams(userID string, characterID pgtype.UUID, categoryID pgtype.UUID, name string) db.CreateCharacterSkillParams {
	return db.CreateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		Name:        name,
		CategoryID:  categoryID,
		BaseValue:   10,
		Value:       35,
		Checked:     false,
		Specialized: false,
		SpecialtyID: pgtype.UUID{},
	}
}

func testUpdateSkillParams(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID, name string) db.UpdateCharacterSkillParams {
	return db.UpdateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		SkillID:     skillID,
		Name:        name,
		CategoryID:  categoryID,
		BaseValue:   15,
		Value:       45,
		Checked:     true,
		Specialized: false,
		SpecialtyID: pgtype.UUID{},
	}
}

func createFinanceTestCreditRatingSkill(t *testing.T, subject *characterIntegrationSubject, userID string, characterID pgtype.UUID) db.CreateCharacterSkillRow {
	t.Helper()

	return createFinanceTestSkill(t, subject, userID, characterID, "Credit Rating")
}

func createFinanceTestSkill(t *testing.T, subject *characterIntegrationSubject, userID string, characterID pgtype.UUID, name string) db.CreateCharacterSkillRow {
	t.Helper()

	categoryID, _ := createSkillTestCategory(t, subject, "Credit Rating")
	skill, err := subject.queries.CreateCharacterSkill(context.Background(), db.CreateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		Name:        name,
		CategoryID:  categoryID,
		BaseValue:   0,
		Value:       35,
		Checked:     false,
		Specialized: false,
		SpecialtyID: pgtype.UUID{},
	})
	require.NoError(t, err)

	return skill
}

func financeString(value string) *string {
	return &value
}

func requireFinanceString(t *testing.T, actual *string, expected string) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func createBackstoryTestBackstory(t *testing.T, subject *characterIntegrationSubject, userID string, characterID pgtype.UUID) db.Backstory {
	t.Helper()

	backstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              userID,
		CharacterID:         characterID,
		PersonalDescription: backstoryString("Test backstory"),
	})
	require.NoError(t, err)

	return backstory
}

func backstoryString(value string) *string {
	return &value
}

func requireBackstoryString(t *testing.T, actual *string, expected string) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}
