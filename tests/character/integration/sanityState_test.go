package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestSanityStateTableUpsertCreatesGetsAndPartiallyUpdatesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	maxSanity := int16(80)
	currentSanity := int16(65)
	tempInsanity := true
	indefInsanity := true

	createdState, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		MaxSanity:     &maxSanity,
		CurrentSanity: &currentSanity,
		TempInsanity:  &tempInsanity,
		IndefInsanity: &indefInsanity,
	})
	require.NoError(t, err)

	require.True(t, createdState.ID.Valid)
	require.Equal(t, character.ID, createdState.CharacterID)
	require.Equal(t, maxSanity, createdState.MaxSanity)
	require.Equal(t, currentSanity, createdState.CurrentSanity)
	require.True(t, createdState.TempInsanity)
	require.True(t, createdState.IndefInsanity)

	fetchedState, err := subject.queries.GetSanityState(context.Background(), db.GetSanityStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, fetchedState.ID)

	falseValue := false
	updates := []struct {
		name string
		set  func(*db.UpsertSanityStateParams)
		want struct {
			temp  bool
			indef bool
		}
	}{
		{
			name: "temporary insanity",
			set:  func(input *db.UpsertSanityStateParams) { input.TempInsanity = &falseValue },
			want: struct{ temp, indef bool }{false, true},
		},
		{
			name: "indefinite insanity",
			set:  func(input *db.UpsertSanityStateParams) { input.IndefInsanity = &falseValue },
			want: struct{ temp, indef bool }{false, false},
		},
	}

	for _, update := range updates {
		t.Run(update.name, func(t *testing.T) {
			input := db.UpsertSanityStateParams{UserID: testUser.ID, CharacterID: character.ID}
			update.set(&input)

			updatedState, err := subject.queries.UpsertSanityState(context.Background(), input)
			require.NoError(t, err)
			require.Equal(t, createdState.ID, updatedState.ID)
			require.Equal(t, maxSanity, updatedState.MaxSanity)
			require.Equal(t, currentSanity, updatedState.CurrentSanity)
			require.Equal(t, update.want.temp, updatedState.TempInsanity)
			require.Equal(t, update.want.indef, updatedState.IndefInsanity)
		})
	}
}

func TestSanityStateTableUpsertUsesDatabaseDefaults(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	state, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, int16(1), state.MaxSanity)
	require.Equal(t, int16(1), state.CurrentSanity)
	require.False(t, state.TempInsanity)
	require.False(t, state.IndefInsanity)
}

func TestSanityStateTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	maxSanity := int16(80)
	_, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		MaxSanity:   &maxSanity,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdState, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		MaxSanity:   &maxSanity,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetSanityState(context.Background(), db.GetSanityStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteSanityState(context.Background(), db.DeleteSanityStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedState, err := subject.queries.DeleteSanityState(context.Background(), db.DeleteSanityStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, deletedState.ID)
}

func TestSanityStateTableRejectsNegativeValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	negative := int16(-1)
	_, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		MaxSanity:   &negative,
	})

	requirePostgresErrorCode(t, err, "23514")
}

func TestSanityStateTableDeletingCharacterCascadesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertSanityState(context.Background(), db.UpsertSanityStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetSanityState(context.Background(), db.GetSanityStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
