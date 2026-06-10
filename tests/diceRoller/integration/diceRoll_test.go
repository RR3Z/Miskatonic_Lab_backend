package tests

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestDiceRollTableCreateAndGet(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	input := createTestDiceRollParams(testUser.ID, character.ID)
	created, err := subject.queries.CreateDiceRoll(context.Background(), input)
	require.NoError(t, err)
	require.True(t, created.ID.Valid)
	require.Equal(t, input.UserID, created.UserID)
	require.Equal(t, input.CharacterID.Bytes, created.CharacterID.Bytes)
	require.Equal(t, input.Expression, created.Expression)
	require.Equal(t, input.Result, created.Result)

	var expectedDetails, actualDetails []map[string]any
	require.NoError(t, json.Unmarshal(input.Details, &expectedDetails))
	require.NoError(t, json.Unmarshal(created.Details, &actualDetails))
	require.Equal(t, expectedDetails, actualDetails)

	require.True(t, created.CreatedAt.Valid)

	fetched, err := subject.queries.GetDiceRoll(context.Background(), db.GetDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		RollID:      created.ID,
	})
	require.NoError(t, err)
	require.Equal(t, created.ID.Bytes, fetched.ID.Bytes)
	require.Equal(t, created.Expression, fetched.Expression)
	require.Equal(t, created.Result, fetched.Result)
}

func TestDiceRollTableCreateEnforcesOwnership(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	ownerUser := createDiceTestUser(t, subject)
	otherUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, ownerUser.ID)

	_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Expression:  "1d20",
		Result:      15,
		Details:     []byte(`[{"type":"dice","sides":20,"rolls":[15]}]`),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDiceRollTableCreateNonexistentCharacter(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)

	_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: diceTestUUID("00000000-0000-0000-0000-000000000000"),
		Expression:  "1d20",
		Result:      15,
		Details:     []byte(`[{"type":"dice","sides":20,"rolls":[15]}]`),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDiceRollTableListReturnsRollsInReverseChronologicalOrder(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	for i := 1; i <= 3; i++ {
		_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
			UserID:      testUser.ID,
			CharacterID: character.ID,
			Expression:  "1d6",
			Result:      int32(i),
			Details:     []byte(`[{"type":"dice","sides":6,"rolls":[` + strconv.Itoa(i) + `]}]`),
		})
		require.NoError(t, err)
	}

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, rolls, 3)
	require.Equal(t, int32(3), rolls[0].Result)
	require.Equal(t, int32(2), rolls[1].Result)
	require.Equal(t, int32(1), rolls[2].Result)
}

func TestDiceRollTableListReturnsEmptyForCharacterWithoutRolls(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, rolls)
}

func TestDiceRollTableListOnlyReturnsRollsForRequestedCharacter(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	char1 := createDiceTestCharacter(t, subject, testUser.ID)
	char2 := createDiceTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, char1.ID))
	require.NoError(t, err)

	_, err = subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, char2.ID))
	require.NoError(t, err)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: char1.ID,
	})
	require.NoError(t, err)
	require.Len(t, rolls, 1)
	require.Equal(t, char1.ID.Bytes, rolls[0].CharacterID.Bytes)
}

func TestDiceRollTableListEnforcesOwnership(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	ownerUser := createDiceTestUser(t, subject)
	otherUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, ownerUser.ID)

	_, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(ownerUser.ID, character.ID))
	require.NoError(t, err)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, rolls)
}

func TestDiceRollTableListLimitsToFiftyRolls(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	for i := 0; i < 55; i++ {
		_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
			UserID:      testUser.ID,
			CharacterID: character.ID,
			Expression:  "1d6",
			Result:      int32(i),
			Details:     []byte(`[{"type":"dice","sides":6,"rolls":[1]}]`),
		})
		require.NoError(t, err)
	}

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, rolls, 50)
}

func TestDiceRollTableGetSingleRoll(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	created, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, character.ID))
	require.NoError(t, err)

	fetched, err := subject.queries.GetDiceRoll(context.Background(), db.GetDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		RollID:      created.ID,
	})
	require.NoError(t, err)
	require.Equal(t, created.ID.Bytes, fetched.ID.Bytes)
}

func TestDiceRollTableGetSingleRollNonexistent(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.GetDiceRoll(context.Background(), db.GetDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		RollID:      diceTestUUID("00000000-0000-0000-0000-000000000000"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDiceRollTableDeleteRoll(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	created, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, character.ID))
	require.NoError(t, err)

	deleted, err := subject.queries.DeleteDiceRoll(context.Background(), db.DeleteDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		RollID:      created.ID,
	})
	require.NoError(t, err)
	require.Equal(t, created.ID.Bytes, deleted.ID.Bytes)

	_, err = subject.queries.GetDiceRoll(context.Background(), db.GetDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		RollID:      created.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDiceRollTableDeleteEnforcesOwnership(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	ownerUser := createDiceTestUser(t, subject)
	otherUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, ownerUser.ID)

	created, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(ownerUser.ID, character.ID))
	require.NoError(t, err)

	_, err = subject.queries.DeleteDiceRoll(context.Background(), db.DeleteDiceRollParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		RollID:      created.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDiceRollTableCleanupKeepsLastFifty(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	for i := 0; i < 55; i++ {
		_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
			UserID:      testUser.ID,
			CharacterID: character.ID,
			Expression:  "1d6",
			Result:      int32(i),
			Details:     []byte(`[{"type":"dice","sides":6,"rolls":[1]}]`),
		})
		require.NoError(t, err)
	}

	err := subject.queries.CleanOldDiceRolls(context.Background(), db.CleanOldDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, rolls, 50)
	for _, r := range rolls {
		require.GreaterOrEqual(t, r.Result, int32(5))
	}
}

func TestDiceRollTableCleanupDoesNothingWhenFiftyOrFewer(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	for i := 0; i < 3; i++ {
		_, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, character.ID))
		require.NoError(t, err)
	}

	err := subject.queries.CleanOldDiceRolls(context.Background(), db.CleanOldDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, rolls, 3)
}

func TestDiceRollTableCleanupEnforcesOwnership(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	ownerUser := createDiceTestUser(t, subject)
	otherUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, ownerUser.ID)

	for i := 0; i < 55; i++ {
		_, err := subject.queries.CreateDiceRoll(context.Background(), db.CreateDiceRollParams{
			UserID:      ownerUser.ID,
			CharacterID: character.ID,
			Expression:  "1d6",
			Result:      int32(i),
			Details:     []byte(`[{"type":"dice","sides":6,"rolls":[1]}]`),
		})
		require.NoError(t, err)
	}

	err := subject.queries.CleanOldDiceRolls(context.Background(), db.CleanOldDiceRollsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	var count int
	err = subject.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM dice_rolls WHERE character_id = $1", character.ID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 55, count)
}

func TestDiceRollTableDeletingCharacterCascadesRolls(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateDiceRoll(context.Background(), createTestDiceRollParams(testUser.ID, character.ID))
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	rolls, err := subject.queries.GetDiceRolls(context.Background(), db.GetDiceRollsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, rolls)
}

func TestDiceRollTableDetailsStoresAndRetrievesJSON(t *testing.T) {
	subject := newDiceIntegrationSubject(t)
	testUser := createDiceTestUser(t, subject)
	character := createDiceTestCharacter(t, subject, testUser.ID)

	input := db.CreateDiceRollParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Expression:  "2d8+1d4",
		Result:      14,
		Details:     []byte(`[{"type":"dice","sides":8,"rolls":[5,7]},{"type":"dice","sides":4,"rolls":[2]}]`),
	}

	created, err := subject.queries.CreateDiceRoll(context.Background(), input)
	require.NoError(t, err)

	var details []map[string]any
	require.NoError(t, json.Unmarshal(created.Details, &details))
	require.Len(t, details, 2)
	require.Equal(t, "dice", details[0]["type"])
	require.Equal(t, float64(8), details[0]["sides"])
	require.Equal(t, "dice", details[1]["type"])
	require.Equal(t, float64(4), details[1]["sides"])
}
