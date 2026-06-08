package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCharacterTableCreateAndGetCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)

	input := testCreateCharacterParams(testUser.ID)

	createdCharacter, err := subject.queries.CreateCharacter(context.Background(), input)
	require.NoError(t, err)

	require.True(t, createdCharacter.ID.Valid)
	require.Equal(t, input.UserID, createdCharacter.UserID)
	require.Equal(t, input.Name, createdCharacter.Name)
	require.Equal(t, input.PlayerName, createdCharacter.PlayerName)
	require.Equal(t, input.Occupation, createdCharacter.Occupation)
	require.Equal(t, input.Age, createdCharacter.Age)
	require.Equal(t, input.Sex, createdCharacter.Sex)
	require.Equal(t, input.Residence, createdCharacter.Residence)
	require.Equal(t, input.Birthplace, createdCharacter.Birthplace)
	require.True(t, createdCharacter.CreatedAt.Valid)
	require.True(t, createdCharacter.UpdatedAt.Valid)

	fetchedCharacter, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: testUser.ID,
		ID:     createdCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacter.ID, fetchedCharacter.ID)
	require.Equal(t, createdCharacter.Name, fetchedCharacter.Name)
}

func TestCharacterTableCreateAllowsNilOptionalFields(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)

	createdCharacter, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{
		UserID: testUser.ID,
		Name:   "Nameless Investigator",
	})
	require.NoError(t, err)

	require.True(t, createdCharacter.ID.Valid)
	require.Equal(t, "Nameless Investigator", createdCharacter.Name)
	require.Nil(t, createdCharacter.PlayerName)
	require.Nil(t, createdCharacter.Occupation)
	require.Nil(t, createdCharacter.Age)
	require.Nil(t, createdCharacter.Sex)
	require.Nil(t, createdCharacter.Residence)
	require.Nil(t, createdCharacter.Birthplace)
}

func TestCharacterTableListsOnlyCharactersForRequestedUser(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	firstUser := createCharacterTestUser(t, subject)
	secondUser := createCharacterTestUser(t, subject)

	firstInput := testCreateCharacterParams(firstUser.ID)
	firstInput.Name = "First User Character"
	firstCharacter, err := subject.queries.CreateCharacter(context.Background(), firstInput)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	secondInput := testCreateCharacterParams(secondUser.ID)
	secondInput.Name = "Second User Character"
	_, err = subject.queries.CreateCharacter(context.Background(), secondInput)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	latestFirstInput := testCreateCharacterParams(firstUser.ID)
	latestFirstInput.Name = "Latest First User Character"
	latestFirstCharacter, err := subject.queries.CreateCharacter(context.Background(), latestFirstInput)
	require.NoError(t, err)

	firstUserCharacters, err := subject.queries.GetAllUserCharacters(context.Background(), firstUser.ID)
	require.NoError(t, err)

	require.Len(t, firstUserCharacters, 2)
	require.Equal(t, latestFirstCharacter.ID, firstUserCharacters[0].ID)
	require.Equal(t, firstCharacter.ID, firstUserCharacters[1].ID)
	for _, character := range firstUserCharacters {
		require.Equal(t, firstUser.ID, character.UserID)
		require.NotEqual(t, secondUser.ID, character.UserID)
	}
}

func TestCharacterTableUpdateCharacterRequiresOwner(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)

	createdCharacter, err := subject.queries.CreateCharacter(context.Background(), testCreateCharacterParams(owner.ID))
	require.NoError(t, err)

	updatedName := "Updated Investigator"
	updatedPlayerName := "Updated Player"
	updatedAge := int16(42)

	_, err = subject.queries.UpdateCharacter(context.Background(), db.UpdateCharacterParams{
		UserID:     otherUser.ID,
		ID:         createdCharacter.ID,
		Name:       updatedName,
		PlayerName: &updatedPlayerName,
		Age:        &updatedAge,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	unchangedCharacter, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: owner.ID,
		ID:     createdCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacter.Name, unchangedCharacter.Name)

	updatedCharacter, err := subject.queries.UpdateCharacter(context.Background(), db.UpdateCharacterParams{
		UserID:     owner.ID,
		ID:         createdCharacter.ID,
		Name:       updatedName,
		PlayerName: &updatedPlayerName,
		Age:        &updatedAge,
	})
	require.NoError(t, err)
	require.Equal(t, updatedName, updatedCharacter.Name)
	require.Equal(t, &updatedPlayerName, updatedCharacter.PlayerName)
	require.Equal(t, &updatedAge, updatedCharacter.Age)
}

func TestCharacterTableDeleteCharacterRequiresOwner(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)

	createdCharacter, err := subject.queries.CreateCharacter(context.Background(), testCreateCharacterParams(owner.ID))
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: otherUser.ID,
		ID:     createdCharacter.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	stillExistingCharacter, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: owner.ID,
		ID:     createdCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacter.ID, stillExistingCharacter.ID)

	deletedCharacter, err := subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: owner.ID,
		ID:     createdCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacter.ID, deletedCharacter.ID)

	_, err = subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: owner.ID,
		ID:     createdCharacter.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCharacterTableRejectsCharacterForMissingUser(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)

	input := testCreateCharacterParams("missing_character_integration_user_" + uniqueCharacterIntegrationSuffix())

	_, err := subject.queries.CreateCharacter(context.Background(), input)

	requirePostgresErrorCode(t, err, "23503")
}

func TestCharacterTableRejectsNegativeAge(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)

	input := testCreateCharacterParams(testUser.ID)
	negativeAge := int16(-1)
	input.Age = &negativeAge

	_, err := subject.queries.CreateCharacter(context.Background(), input)

	requirePostgresErrorCode(t, err, "23514")
}

func TestCharacterTableDeletingUserCascadesCharacters(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)

	createdCharacter, err := subject.queries.CreateCharacter(context.Background(), testCreateCharacterParams(testUser.ID))
	require.NoError(t, err)

	err = subject.queries.DeleteUserByClerkID(context.Background(), testUser.ID)
	require.NoError(t, err)

	_, err = subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: testUser.ID,
		ID:     createdCharacter.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
