package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestMagicStateTableUpsertCreatesGetsAndPartiallyUpdatesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	maxMp := int16(18)
	currentMp := int16(11)

	createdState, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		MaxMp:       &maxMp,
		CurrentMp:   &currentMp,
	})
	require.NoError(t, err)

	require.True(t, createdState.ID.Valid)
	require.Equal(t, character.ID, createdState.CharacterID)
	require.Equal(t, maxMp, createdState.MaxMp)
	require.Equal(t, currentMp, createdState.CurrentMp)

	fetchedState, err := subject.queries.GetMagicState(context.Background(), db.GetMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, fetchedState.ID)

	updatedCurrentMp := int16(5)
	updatedState, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		CurrentMp:   &updatedCurrentMp,
	})
	require.NoError(t, err)

	require.Equal(t, createdState.ID, updatedState.ID)
	require.Equal(t, maxMp, updatedState.MaxMp)
	require.Equal(t, updatedCurrentMp, updatedState.CurrentMp)
}

func TestMagicStateTableUpsertUsesDatabaseDefaults(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	state, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, int16(1), state.MaxMp)
	require.Equal(t, int16(1), state.CurrentMp)
}

func TestMagicStateTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	maxMp := int16(18)
	_, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		MaxMp:       &maxMp,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdState, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		MaxMp:       &maxMp,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetMagicState(context.Background(), db.GetMagicStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteMagicState(context.Background(), db.DeleteMagicStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedState, err := subject.queries.DeleteMagicState(context.Background(), db.DeleteMagicStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, deletedState.ID)
}

func TestMagicStateTableRejectsNegativeValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	negative := int16(-1)
	_, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		MaxMp:       &negative,
	})

	requirePostgresErrorCode(t, err, "23514")
}

func TestMagicStateTableDeletingCharacterCascadesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertMagicState(context.Background(), db.UpsertMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetMagicState(context.Background(), db.GetMagicStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
