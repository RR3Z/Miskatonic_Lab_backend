package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestHealthStateTableUpsertCreatesGetsAndPartiallyUpdatesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	maxHp := int16(12)
	currentHp := int16(7)
	majorWound := true
	unconscious := true
	dying := true
	dead := true

	createdState, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		MaxHp:       &maxHp,
		CurrentHp:   &currentHp,
		MajorWound:  &majorWound,
		Unconscious: &unconscious,
		Dying:       &dying,
		Dead:        &dead,
	})
	require.NoError(t, err)

	require.True(t, createdState.ID.Valid)
	require.Equal(t, character.ID, createdState.CharacterID)
	require.Equal(t, maxHp, createdState.MaxHp)
	require.Equal(t, currentHp, createdState.CurrentHp)
	require.True(t, createdState.MajorWound)
	require.True(t, createdState.Unconscious)
	require.True(t, createdState.Dying)
	require.True(t, createdState.Dead)

	fetchedState, err := subject.queries.GetHealthState(context.Background(), db.GetHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, fetchedState.ID)

	falseValue := false
	updates := []struct {
		name string
		set  func(*db.UpsertHealthStateParams)
		want struct {
			majorWound  bool
			unconscious bool
			dying       bool
			dead        bool
		}
	}{
		{
			name: "major wound",
			set:  func(input *db.UpsertHealthStateParams) { input.MajorWound = &falseValue },
			want: struct{ majorWound, unconscious, dying, dead bool }{false, true, true, true},
		},
		{
			name: "unconscious",
			set:  func(input *db.UpsertHealthStateParams) { input.Unconscious = &falseValue },
			want: struct{ majorWound, unconscious, dying, dead bool }{false, false, true, true},
		},
		{
			name: "dying",
			set:  func(input *db.UpsertHealthStateParams) { input.Dying = &falseValue },
			want: struct{ majorWound, unconscious, dying, dead bool }{false, false, false, true},
		},
		{
			name: "dead",
			set:  func(input *db.UpsertHealthStateParams) { input.Dead = &falseValue },
			want: struct{ majorWound, unconscious, dying, dead bool }{false, false, false, false},
		},
	}

	for _, update := range updates {
		t.Run(update.name, func(t *testing.T) {
			input := db.UpsertHealthStateParams{UserID: testUser.ID, CharacterID: character.ID}
			update.set(&input)

			updatedState, err := subject.queries.UpsertHealthState(context.Background(), input)
			require.NoError(t, err)
			require.Equal(t, createdState.ID, updatedState.ID)
			require.Equal(t, maxHp, updatedState.MaxHp)
			require.Equal(t, currentHp, updatedState.CurrentHp)
			require.Equal(t, update.want.majorWound, updatedState.MajorWound)
			require.Equal(t, update.want.unconscious, updatedState.Unconscious)
			require.Equal(t, update.want.dying, updatedState.Dying)
			require.Equal(t, update.want.dead, updatedState.Dead)
		})
	}
}

func TestHealthStateTableUpsertUsesDatabaseDefaults(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	state, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, int16(1), state.MaxHp)
	require.Equal(t, int16(1), state.CurrentHp)
	require.False(t, state.MajorWound)
	require.False(t, state.Unconscious)
	require.False(t, state.Dying)
	require.False(t, state.Dead)
}

func TestHealthStateTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	maxHp := int16(10)
	_, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		MaxHp:       &maxHp,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdState, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		MaxHp:       &maxHp,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetHealthState(context.Background(), db.GetHealthStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteHealthState(context.Background(), db.DeleteHealthStateParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedState, err := subject.queries.DeleteHealthState(context.Background(), db.DeleteHealthStateParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdState.ID, deletedState.ID)
}

func TestHealthStateTableRejectsNegativeValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	negative := int16(-1)
	_, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		MaxHp:       &negative,
	})

	requirePostgresErrorCode(t, err, "23514")
}

func TestHealthStateTableDeletingCharacterCascadesState(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetHealthState(context.Background(), db.GetHealthStateParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
