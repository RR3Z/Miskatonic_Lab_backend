package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCharacteristicsTableUpsertCreatesGetsAndReplacesCharacteristics(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:       testUser.ID,
		CharacterID:  character.ID,
		Strength:     characterInt16(60),
		Constitution: characterInt16(55),
		Size:         characterInt16(70),
		Dexterity:    characterInt16(45),
		Appearance:   characterInt16(50),
		Intelligence: characterInt16(80),
		Power:        characterInt16(65),
		Education:    characterInt16(75),
	})
	require.NoError(t, err)

	require.True(t, createdCharacteristics.ID.Valid)
	require.Equal(t, character.ID, createdCharacteristics.CharacterID)
	requireCharacteristicValue(t, createdCharacteristics.Strength, 60)
	requireCharacteristicValue(t, createdCharacteristics.Constitution, 55)
	requireCharacteristicValue(t, createdCharacteristics.Size, 70)
	requireCharacteristicValue(t, createdCharacteristics.Dexterity, 45)
	requireCharacteristicValue(t, createdCharacteristics.Appearance, 50)
	requireCharacteristicValue(t, createdCharacteristics.Intelligence, 80)
	requireCharacteristicValue(t, createdCharacteristics.Power, 65)
	requireCharacteristicValue(t, createdCharacteristics.Education, 75)
	require.True(t, createdCharacteristics.CreatedAt.Valid)
	require.True(t, createdCharacteristics.UpdatedAt.Valid)

	fetchedCharacteristics, err := subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacteristics.ID, fetchedCharacteristics.ID)

	time.Sleep(5 * time.Millisecond)

	updatedCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(40),
		Power:       characterInt16(90),
	})
	require.NoError(t, err)

	require.Equal(t, createdCharacteristics.ID, updatedCharacteristics.ID)
	requireCharacteristicValue(t, updatedCharacteristics.Strength, 40)
	require.Nil(t, updatedCharacteristics.Constitution)
	require.Nil(t, updatedCharacteristics.Size)
	require.Nil(t, updatedCharacteristics.Dexterity)
	require.Nil(t, updatedCharacteristics.Appearance)
	require.Nil(t, updatedCharacteristics.Intelligence)
	requireCharacteristicValue(t, updatedCharacteristics.Power, 90)
	require.Nil(t, updatedCharacteristics.Education)
	require.True(t, updatedCharacteristics.UpdatedAt.Time.After(createdCharacteristics.UpdatedAt.Time) || updatedCharacteristics.UpdatedAt.Time.Equal(createdCharacteristics.UpdatedAt.Time))
}

func TestCharacteristicsTableUpsertAllowsAllNilValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	characteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.True(t, characteristics.ID.Valid)
	require.Equal(t, character.ID, characteristics.CharacterID)
	require.Nil(t, characteristics.Strength)
	require.Nil(t, characteristics.Constitution)
	require.Nil(t, characteristics.Size)
	require.Nil(t, characteristics.Dexterity)
	require.Nil(t, characteristics.Appearance)
	require.Nil(t, characteristics.Intelligence)
	require.Nil(t, characteristics.Power)
	require.Nil(t, characteristics.Education)
}

func TestCharacteristicsTableAllowsZeroValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	characteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:       testUser.ID,
		CharacterID:  character.ID,
		Strength:     characterInt16(0),
		Constitution: characterInt16(0),
		Size:         characterInt16(0),
		Dexterity:    characterInt16(0),
		Appearance:   characterInt16(0),
		Intelligence: characterInt16(0),
		Power:        characterInt16(0),
		Education:    characterInt16(0),
	})
	require.NoError(t, err)

	requireCharacteristicValue(t, characteristics.Strength, 0)
	requireCharacteristicValue(t, characteristics.Constitution, 0)
	requireCharacteristicValue(t, characteristics.Size, 0)
	requireCharacteristicValue(t, characteristics.Dexterity, 0)
	requireCharacteristicValue(t, characteristics.Appearance, 0)
	requireCharacteristicValue(t, characteristics.Intelligence, 0)
	requireCharacteristicValue(t, characteristics.Power, 0)
	requireCharacteristicValue(t, characteristics.Education, 0)
}

func TestCharacteristicsTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
	})
	require.NoError(t, err)

	_, err = subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacteristics(context.Background(), db.DeleteCharacteristicsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedCharacteristics, err := subject.queries.DeleteCharacteristics(context.Background(), db.DeleteCharacteristicsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacteristics.ID, deletedCharacteristics.ID)
}

func TestCharacteristicsTableUnauthorizedUpsertDoesNotMutateExistingCharacteristics(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	createdCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Power:       characterInt16(70),
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(10),
		Power:       characterInt16(20),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedCharacteristics, err := subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdCharacteristics.ID, fetchedCharacteristics.ID)
	requireCharacteristicValue(t, fetchedCharacteristics.Strength, 60)
	requireCharacteristicValue(t, fetchedCharacteristics.Power, 70)
}

func TestCharacteristicsTableReturnsNoRowsBeforeUpsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacteristics(context.Background(), db.DeleteCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCharacteristicsTableKeepsCharacteristicsScopedToRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
		Strength:    characterInt16(40),
	})
	require.NoError(t, err)

	secondCharacteristics, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
		Strength:    characterInt16(80),
	})
	require.NoError(t, err)

	fetchedFirstCharacteristics, err := subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstCharacteristics.ID, fetchedFirstCharacteristics.ID)
	requireCharacteristicValue(t, fetchedFirstCharacteristics.Strength, 40)

	fetchedSecondCharacteristics, err := subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, secondCharacteristics.ID, fetchedSecondCharacteristics.ID)
	requireCharacteristicValue(t, fetchedSecondCharacteristics.Strength, 80)
}

func TestCharacteristicsTableReturnsNoRowsForMissingCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	missingCharacterID := characterTestUUID("cccccccc-cccc-cccc-cccc-cccccccccccc")

	_, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
		Strength:    characterInt16(60),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacteristics(context.Background(), db.DeleteCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCharacteristicsTableRejectsNegativeValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(-1),
	})
	requirePostgresErrorCode(t, err, "23514")
}

func TestCharacteristicsTableDeletingCharacterCascadesCharacteristics(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertCharacteristics(context.Background(), db.UpsertCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetCharacteristics(context.Background(), db.GetCharacteristicsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func characterInt16(value int16) *int16 {
	return &value
}

func requireCharacteristicValue(t *testing.T, actual *int16, expected int16) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}
