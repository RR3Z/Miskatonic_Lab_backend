package tests

import (
	"testing"

	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newTestCharacterServiceForHealth() (*FakeHealthDBTX, *characterServices.CharacterService) {
	dbtx := &FakeHealthDBTX{QueryRowData: healthRowData(testHealthState())}
	repos := &repository.Repository{Queries: db.New(dbtx)}

	return dbtx, characterServices.NewCharacterService(repos)
}

func testHealthState() db.HealthState {
	return db.HealthState{
		ID:          testHealthUUID("22222222-2222-2222-2222-222222222222"),
		CharacterID: testHealthUUID("11111111-1111-1111-1111-111111111111"),
		MaxHp:       10,
		CurrentHp:   7,
		CreatedAt:   testHealthTimestamptz(),
		UpdatedAt:   testHealthTimestamptz(),
	}
}

func testUpsertHealthInput(maxHp *int16, currentHp *int16) healthDTO.UpsertHealthInput {
	return healthDTO.UpsertHealthInput{
		UserID:      "user_1",
		CharacterID: testHealthUUID("11111111-1111-1111-1111-111111111111"),
		MaxHp:       maxHp,
		CurrentHp:   currentHp,
	}
}

func healthInt16(value int16) *int16 {
	return &value
}

func healthBool(value bool) *bool {
	return &value
}

func testHealthTimestamptz() pgtype.Timestamptz {
	var value pgtype.Timestamptz
	err := value.Scan("2026-06-07 12:00:00+03")
	if err != nil {
		panic(err)
	}

	return value
}

func testHealthUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func requireSameHealthState(t *testing.T, expected db.HealthState, actual db.HealthState) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxHp, actual.MaxHp)
	require.Equal(t, expected.CurrentHp, actual.CurrentHp)
	require.Equal(t, expected.MajorWound, actual.MajorWound)
	require.Equal(t, expected.Unconscious, actual.Unconscious)
	require.Equal(t, expected.Dying, actual.Dying)
	require.Equal(t, expected.Dead, actual.Dead)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}

func healthRowData(health db.HealthState) []any {
	return []any{
		health.ID,
		health.CharacterID,
		health.MaxHp,
		health.CurrentHp,
		health.MajorWound,
		health.Unconscious,
		health.Dying,
		health.Dead,
		health.CreatedAt,
		health.UpdatedAt,
	}
}
