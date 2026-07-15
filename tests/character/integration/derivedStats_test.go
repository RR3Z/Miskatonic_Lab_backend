package tests

import (
	"context"
	"testing"
	"time"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestDerivedStatsTableUpsertCreatesGetsAndPartiallyUpdatesStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
		Physique:    characterInt16(1),
		DamageBonus: characterString("+1d6"),
		DodgeValue:  characterInt16(40),
	})
	require.NoError(t, err)

	require.True(t, createdStats.ID.Valid)
	require.Equal(t, character.ID, createdStats.CharacterID)
	requireDerivedStatValue(t, createdStats.Speed, 8)
	requireDerivedStatValue(t, createdStats.Physique, 1)
	requireDerivedStatString(t, createdStats.DamageBonus, "+1d6")
	requireDerivedStatValue(t, createdStats.DodgeValue, 40)
	require.True(t, createdStats.CreatedAt.Valid)
	require.True(t, createdStats.UpdatedAt.Valid)

	fetchedStats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdStats.ID, fetchedStats.ID)

	time.Sleep(5 * time.Millisecond)

	updatedStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(9),
		DodgeValue:  characterInt16(45),
	})
	require.NoError(t, err)

	require.Equal(t, createdStats.ID, updatedStats.ID)
	requireDerivedStatValue(t, updatedStats.Speed, 9)
	requireDerivedStatValue(t, updatedStats.Physique, 1)
	requireDerivedStatString(t, updatedStats.DamageBonus, "+1d6")
	requireDerivedStatValue(t, updatedStats.DodgeValue, 45)
	require.True(t, updatedStats.UpdatedAt.Time.After(createdStats.UpdatedAt.Time) || updatedStats.UpdatedAt.Time.Equal(createdStats.UpdatedAt.Time))
}

func TestDerivedStatsTableUpsertAllowsAllNilValuesOnInsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	stats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.True(t, stats.ID.Valid)
	require.Equal(t, character.ID, stats.CharacterID)
	require.Nil(t, stats.Speed)
	require.Nil(t, stats.Physique)
	require.Nil(t, stats.DamageBonus)
	require.Nil(t, stats.DodgeValue)
}

func TestDerivedStatsTableNilUpdateDoesNotOverwriteExistingValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
		Physique:    characterInt16(1),
		DamageBonus: characterString("+1d6"),
		DodgeValue:  characterInt16(40),
	})
	require.NoError(t, err)

	updatedStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, createdStats.ID, updatedStats.ID)
	requireDerivedStatValue(t, updatedStats.Speed, 8)
	requireDerivedStatValue(t, updatedStats.Physique, 1)
	requireDerivedStatString(t, updatedStats.DamageBonus, "+1d6")
	requireDerivedStatValue(t, updatedStats.DodgeValue, 40)
}

func TestDerivedStatsTablePartialUpdateAfterNilInsertOnlySetsProvidedValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Nil(t, createdStats.Speed)
	require.Nil(t, createdStats.Physique)
	require.Nil(t, createdStats.DamageBonus)
	require.Nil(t, createdStats.DodgeValue)

	updatedStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Physique:    characterInt16(1),
	})
	require.NoError(t, err)

	require.Equal(t, createdStats.ID, updatedStats.ID)
	require.Nil(t, updatedStats.Speed)
	requireDerivedStatValue(t, updatedStats.Physique, 1)
	require.Nil(t, updatedStats.DamageBonus)
	require.Nil(t, updatedStats.DodgeValue)
}

func TestDerivedStatsTableAllowsZeroValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	stats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(0),
		Physique:    characterInt16(0),
		DamageBonus: characterString("0"),
		DodgeValue:  characterInt16(0),
	})
	require.NoError(t, err)

	requireDerivedStatValue(t, stats.Speed, 0)
	requireDerivedStatValue(t, stats.Physique, 0)
	requireDerivedStatString(t, stats.DamageBonus, "0")
	requireDerivedStatValue(t, stats.DodgeValue, 0)
}

func TestDerivedStatsTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
	})
	require.NoError(t, err)

	_, err = subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteDerivedStats(context.Background(), db.DeleteDerivedStatsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedStats, err := subject.queries.DeleteDerivedStats(context.Background(), db.DeleteDerivedStatsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdStats.ID, deletedStats.ID)
}

func TestDerivedStatsTableUnauthorizedUpsertDoesNotMutateExistingStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
		DodgeValue:  characterInt16(40),
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(1),
		DodgeValue:  characterInt16(5),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedStats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdStats.ID, fetchedStats.ID)
	requireDerivedStatValue(t, fetchedStats.Speed, 8)
	requireDerivedStatValue(t, fetchedStats.DodgeValue, 40)
}

func TestDerivedStatsTableReturnsNoRowsBeforeUpsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteDerivedStats(context.Background(), db.DeleteDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDerivedStatsTableKeepsStatsScopedToRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
		Speed:       characterInt16(8),
	})
	require.NoError(t, err)

	secondStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
		Speed:       characterInt16(10),
	})
	require.NoError(t, err)

	fetchedFirstStats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstStats.ID, fetchedFirstStats.ID)
	requireDerivedStatValue(t, fetchedFirstStats.Speed, 8)

	fetchedSecondStats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, secondStats.ID, fetchedSecondStats.ID)
	requireDerivedStatValue(t, fetchedSecondStats.Speed, 10)
}

func TestDerivedStatsTableReturnsNoRowsForMissingCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	missingCharacterID := characterTestUUID("dddddddd-dddd-dddd-dddd-dddddddddddd")

	_, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
		Speed:       characterInt16(8),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteDerivedStats(context.Background(), db.DeleteDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDerivedStatsTableDeleteReturnsDeletedValuesAndAllowsRecreate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
		DodgeValue:  characterInt16(40),
	})
	require.NoError(t, err)

	deletedStats, err := subject.queries.DeleteDerivedStats(context.Background(), db.DeleteDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, createdStats.ID, deletedStats.ID)
	requireDerivedStatValue(t, deletedStats.Speed, 8)
	requireDerivedStatValue(t, deletedStats.DodgeValue, 40)

	recreatedStats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(10),
	})
	require.NoError(t, err)

	require.NotEqual(t, deletedStats.ID, recreatedStats.ID)
	requireDerivedStatValue(t, recreatedStats.Speed, 10)
}

func TestDerivedStatsTableAllowsNegativePhysiqueAndDamageBonusValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	stats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Physique:    characterInt16(-2),
		DamageBonus: characterString("-2"),
	})
	require.NoError(t, err)

	requireDerivedStatValue(t, stats.Physique, -2)
	requireDerivedStatString(t, stats.DamageBonus, "-2")
}

func TestDerivedStatsTableAllowsDiceDamageBonusValues(t *testing.T) {
	tests := []string{"+1d4", "+1d6", "+2d6", "+10d6"}

	for _, damageBonus := range tests {
		t.Run(damageBonus, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)

			stats, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
				UserID:      testUser.ID,
				CharacterID: character.ID,
				DamageBonus: characterString(damageBonus),
			})
			require.NoError(t, err)
			requireDerivedStatString(t, stats.DamageBonus, damageBonus)
		})
	}
}

func TestDerivedStatsTableRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name   string
		params func(userID string, character db.Character) db.UpsertDerivedStatsParams
	}{
		{
			name: "speed",
			params: func(userID string, character db.Character) db.UpsertDerivedStatsParams {
				return db.UpsertDerivedStatsParams{
					UserID:      userID,
					CharacterID: character.ID,
					Speed:       characterInt16(-1),
				}
			},
		},
		{
			name: "damage bonus",
			params: func(userID string, character db.Character) db.UpsertDerivedStatsParams {
				return db.UpsertDerivedStatsParams{
					UserID:      userID,
					CharacterID: character.ID,
					DamageBonus: characterString("not-a-damage-bonus"),
				}
			},
		},
		{
			name: "dodge value",
			params: func(userID string, character db.Character) db.UpsertDerivedStatsParams {
				return db.UpsertDerivedStatsParams{
					UserID:      userID,
					CharacterID: character.ID,
					DodgeValue:  characterInt16(-1),
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)

			_, err := subject.queries.UpsertDerivedStats(context.Background(), tc.params(testUser.ID, character))
			requirePostgresErrorCode(t, err, "23514")
		})
	}
}

func TestDerivedStatsTableDeletingCharacterCascadesStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertDerivedStats(context.Background(), db.UpsertDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Speed:       characterInt16(8),
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCharacterServiceUpsertCharacteristicsRecalculatesDerivedStatsWithoutAge(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{
		UserID: testUser.ID,
		Name:   "Ageless Investigator",
	})
	require.NoError(t, err)

	recorder := &characterEventRecorder{}
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, recorder)

	characteristics, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Size:        characterInt16(40),
		Dexterity:   characterInt16(70),
	})
	require.NoError(t, err)
	requireCharacteristicValue(t, characteristics.Strength, 60)

	stats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireDerivedStatValue(t, stats.Speed, 9)
	requireDerivedStatValue(t, stats.Physique, 0)
	requireDerivedStatString(t, stats.DamageBonus, "0")
	requireDerivedStatValue(t, stats.DodgeValue, 35)

	event := requireLastCharacterEvent[characterEvents.CharacterDerivedStatsAutoRecalculateSucceeded](t, recorder)
	require.Equal(t, testUser.ID, event.UserID)
	require.Equal(t, character.ID.String(), event.CharacterID)
	require.Equal(t, "characteristics_upsert", event.Source)
}

func TestCharacterServiceUpsertCharacteristicsSkipsDerivedStatsCalculationWithoutFormulaCharacteristics(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	recorder := &characterEventRecorder{}
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, recorder)

	characteristics, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Size:        characterInt16(40),
	})
	require.NoError(t, err)
	requireCharacteristicValue(t, characteristics.Strength, 60)
	require.Nil(t, characteristics.Dexterity)

	_, err = subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	event := requireLastCharacterEvent[characterEvents.CharacterDerivedStatsAutoRecalculateSkipped](t, recorder)
	require.Equal(t, testUser.ID, event.UserID)
	require.Equal(t, character.ID.String(), event.CharacterID)
	require.Equal(t, "characteristics_upsert", event.Source)
	require.Equal(t, "required_characteristics_missing", event.Reason)
}

func TestCharacterServiceUpsertCharacteristicsRecalculatesDerivedStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	recorder := &characterEventRecorder{}
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, recorder)

	_, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Size:        characterInt16(40),
		Dexterity:   characterInt16(70),
	})
	require.NoError(t, err)

	stats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireDerivedStatValue(t, stats.Speed, 9)
	requireDerivedStatValue(t, stats.Physique, 0)
	requireDerivedStatString(t, stats.DamageBonus, "0")
	requireDerivedStatValue(t, stats.DodgeValue, 35)

	event := requireLastCharacterEvent[characterEvents.CharacterDerivedStatsAutoRecalculateSucceeded](t, recorder)
	require.Equal(t, testUser.ID, event.UserID)
	require.Equal(t, character.ID.String(), event.CharacterID)
	require.Equal(t, "characteristics_upsert", event.Source)
}

func TestCharacterServiceUpdateCharacterAgeDoesNotRecalculateDerivedStats(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	recorder := &characterEventRecorder{}
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, recorder)

	_, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Size:        characterInt16(40),
		Dexterity:   characterInt16(70),
	})
	require.NoError(t, err)

	recorder.events = nil

	age := int16(80)
	_, err = service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID:     testUser.ID,
		ID:         character.ID,
		Name:       character.Name,
		PlayerName: character.PlayerName,
		Occupation: character.Occupation,
		Age:        &age,
		Sex:        character.Sex,
		Residence:  character.Residence,
		Birthplace: character.Birthplace,
	})
	require.NoError(t, err)

	stats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireDerivedStatValue(t, stats.Speed, 9)
	require.Empty(t, recorder.events)
}

func TestCharacterServiceUpdateCharacterDoesNotRecalculateDerivedStatsWhenAgeCleared(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	recorder := &characterEventRecorder{}
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, recorder)

	_, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Strength:    characterInt16(60),
		Size:        characterInt16(40),
		Dexterity:   characterInt16(70),
	})
	require.NoError(t, err)

	recorder.events = nil

	_, err = service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID:     testUser.ID,
		ID:         character.ID,
		Name:       character.Name,
		PlayerName: character.PlayerName,
		Occupation: character.Occupation,
		Age:        nil,
		Sex:        character.Sex,
		Residence:  character.Residence,
		Birthplace: character.Birthplace,
	})
	require.NoError(t, err)

	updatedCharacter, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)
	require.Nil(t, updatedCharacter.Age)

	stats, err := subject.queries.GetDerivedStats(context.Background(), db.GetDerivedStatsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireDerivedStatValue(t, stats.Speed, 9)

	require.Empty(t, recorder.events)
}
