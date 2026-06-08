package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestLuckStateTableUpsertCreatesGetsAndPartiallyUpdatesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	startingLuck := int16(70)
	currentLuck := int16(55)

	createdState, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:       testUser.ID,
		CharacterID:  character.ID,
		StartingLuck: &startingLuck,
		CurrentLuck:  &currentLuck,
	})
	require.NoError(t, err)

	require.True(t, createdState.ID.Valid)
	require.Equal(t, character.ID, createdState.CharacterID)
	require.Equal(t, startingLuck, createdState.StartingLuck)
	require.Equal(t, currentLuck, createdState.CurrentLuck)

	fetchedState, err := subject.queries.GetLuckState(context.Background(), db.GetLuckStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, fetchedState.ID)

	updatedCurrentLuck := int16(40)
	updatedState, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		CurrentLuck: &updatedCurrentLuck,
	})
	require.NoError(t, err)

	require.Equal(t, createdState.ID, updatedState.ID)
	require.Equal(t, startingLuck, updatedState.StartingLuck)
	require.Equal(t, updatedCurrentLuck, updatedState.CurrentLuck)
}

func TestLuckStateTableUpsertUsesDatabaseDefaults(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	state, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, int16(1), state.StartingLuck)
	require.Equal(t, int16(1), state.CurrentLuck)
}

func TestLuckStateTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	startingLuck := int16(70)
	_, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:       otherUser.ID,
		CharacterID:  character.ID,
		StartingLuck: &startingLuck,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdState, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:       owner.ID,
		CharacterID:  character.ID,
		StartingLuck: &startingLuck,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetLuckState(context.Background(), db.GetLuckStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteLuckState(context.Background(), db.DeleteLuckStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedState, err := subject.queries.DeleteLuckState(context.Background(), db.DeleteLuckStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, deletedState.ID)
}

func TestLuckStateTableRejectsNegativeValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	negative := int16(-1)
	_, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:       testUser.ID,
		CharacterID:  character.ID,
		StartingLuck: &negative,
	})

	requirePostgresErrorCode(t, err, "23514")
}

func TestLuckStateTableDeletingCharacterCascadesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertLuckState(context.Background(), db.UpsertLuckStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetLuckState(context.Background(), db.GetLuckStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
