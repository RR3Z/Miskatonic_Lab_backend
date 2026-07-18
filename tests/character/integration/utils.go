package tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool := testdb.Open(t)

	repos := repository.NewRepository(pool)
	return &characterIntegrationSubject{
		pool:    pool,
		queries: repos.Queries,
	}
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
	occupation := "Antiquarian"
	age := int16(37)
	sex := "male"
	residence := "Arkham"
	birthplace := "Boston"
	return db.CreateCharacterParams{
		UserID:     userID,
		Name:       "Dr. Henry Armitage",
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

func testCreateSkillParams(userID string, characterID pgtype.UUID, name string) db.CreateCharacterSkillParams {
	return db.CreateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		Name:        name,
		BaseValue:   10,
		Value:       35,
		Checked:     false,
	}
}

func testUpdateSkillParams(userID string, characterID pgtype.UUID, skillID pgtype.UUID, name string) db.UpdateCharacterSkillParams {
	return db.UpdateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		SkillID:     skillID,
		Name:        name,
		BaseValue:   15,
		Value:       45,
		Checked:     true,
	}
}

func createFinanceTestSkill(t *testing.T, subject *characterIntegrationSubject, userID string, characterID pgtype.UUID, name string) db.CreateCharacterSkillRow {
	t.Helper()

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), db.CreateCharacterSkillParams{
		UserID:      userID,
		CharacterID: characterID,
		Name:        name,
		BaseValue:   0,
		Value:       35,
		Checked:     false,
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
