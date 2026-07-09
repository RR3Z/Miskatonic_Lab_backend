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
	require.Equal(t, input.PortraitUrl, createdCharacter.PortraitUrl)
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
	require.Nil(t, createdCharacter.PortraitUrl)
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

func TestCharacterTableListsCharacterCardsWithPortraitAndStateStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)

	noStatesInput := testCreateCharacterParams(testUser.ID)
	noStatesInput.Name = "No State Rows"
	noStatesInput.PortraitUrl = nil
	emptySex := ""
	noStatesInput.Sex = &emptySex
	noStatesCharacter, err := subject.queries.CreateCharacter(context.Background(), noStatesInput)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	portraitURL := "https://assets.example.test/portraits/card.webp"
	withStatesInput := testCreateCharacterParams(testUser.ID)
	withStatesInput.Name = "Card Ready"
	withStatesInput.PortraitUrl = &portraitURL
	withStatesCharacter, err := subject.queries.CreateCharacter(context.Background(), withStatesInput)
	require.NoError(t, err)

	_, err = subject.queries.CreateCharacter(context.Background(), testCreateCharacterParams(otherUser.ID))
	require.NoError(t, err)

	_, err = subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: withStatesCharacter.ID,
		MaxHp:       characterInt16(12),
		CurrentHp:   characterInt16(7),
	})
	require.NoError(t, err)
	_, err = subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: withStatesCharacter.ID,
		MaxMp:       characterInt16(9),
		CurrentMp:   characterInt16(4),
	})
	require.NoError(t, err)
	_, err = subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:        testUser.ID,
		CharacterID:   withStatesCharacter.ID,
		MaxSanity:     characterInt16(60),
		CurrentSanity: characterInt16(33),
	})
	require.NoError(t, err)
	_, err = subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:       testUser.ID,
		CharacterID:  withStatesCharacter.ID,
		StartingLuck: characterInt16(45),
		CurrentLuck:  characterInt16(20),
	})
	require.NoError(t, err)

	cards, err := subject.queries.GetAllUserCharacterCards(context.Background(), testUser.ID)
	require.NoError(t, err)

	require.Len(t, cards, 2)
	require.Equal(t, withStatesCharacter.ID, cards[0].ID)
	require.Equal(t, "Card Ready", cards[0].Name)
	require.Equal(t, withStatesInput.Occupation, cards[0].Occupation)
	require.Equal(t, withStatesInput.Age, cards[0].Age)
	require.Equal(t, withStatesInput.Sex, cards[0].Sex)
	require.Equal(t, withStatesInput.Residence, cards[0].Residence)
	require.Equal(t, &portraitURL, cards[0].PortraitUrl)
	require.Equal(t, int16(7), cards[0].CurrentHp)
	require.Equal(t, int16(12), cards[0].MaxHp)
	require.Equal(t, int16(4), cards[0].CurrentMp)
	require.Equal(t, int16(9), cards[0].MaxMp)
	require.Equal(t, int16(33), cards[0].CurrentSanity)
	require.Equal(t, int16(60), cards[0].MaxSanity)
	require.Equal(t, int16(20), cards[0].CurrentLuck)
	require.Equal(t, int16(45), cards[0].StartingLuck)

	require.Equal(t, noStatesCharacter.ID, cards[1].ID)
	require.Equal(t, "No State Rows", cards[1].Name)
	require.Equal(t, &emptySex, cards[1].Sex)
	require.Nil(t, cards[1].PortraitUrl)
	require.Equal(t, int16(0), cards[1].CurrentHp)
	require.Equal(t, int16(0), cards[1].MaxHp)
	require.Equal(t, int16(0), cards[1].CurrentMp)
	require.Equal(t, int16(0), cards[1].MaxMp)
	require.Equal(t, int16(0), cards[1].CurrentSanity)
	require.Equal(t, int16(0), cards[1].MaxSanity)
	require.Equal(t, int16(0), cards[1].CurrentLuck)
	require.Equal(t, int16(0), cards[1].StartingLuck)
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
