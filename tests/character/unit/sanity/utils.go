package tests

import (
	"testing"

	characterModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newTestCharacterServiceForSanity() (*FakeSanityDBTX, *characterServices.CharacterService) {
	dbtx := &FakeSanityDBTX{QueryRowData: sanityRowData(testSanityState())}
	repos := &repository.Repository{Queries: db.New(dbtx)}

	return dbtx, characterServices.NewCharacterService(repos)
}

func testSanityState() db.SanityState {
	return db.SanityState{
		ID:            testSanityUUID("22222222-2222-2222-2222-222222222222"),
		CharacterID:   testSanityUUID("11111111-1111-1111-1111-111111111111"),
		MaxSanity:     60,
		CurrentSanity: 40,
		CreatedAt:     testSanityTimestamptz(),
		UpdatedAt:     testSanityTimestamptz(),
	}
}

func testUpsertSanityInput(maxSanity *int16, currentSanity *int16) characterModel.UpsertSanityInput {
	return characterModel.UpsertSanityInput{
		UserID:        "user_1",
		CharacterID:   testSanityUUID("11111111-1111-1111-1111-111111111111"),
		MaxSanity:     maxSanity,
		CurrentSanity: currentSanity,
	}
}

func sanityInt16(value int16) *int16 {
	return &value
}

func sanityBool(value bool) *bool {
	return &value
}

func testSanityTimestamptz() pgtype.Timestamptz {
	var value pgtype.Timestamptz
	err := value.Scan("2026-06-07 12:00:00+03")
	if err != nil {
		panic(err)
	}

	return value
}

func testSanityUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func requireSameSanityState(t *testing.T, expected db.SanityState, actual db.SanityState) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxSanity, actual.MaxSanity)
	require.Equal(t, expected.CurrentSanity, actual.CurrentSanity)
	require.Equal(t, expected.TempInsanity, actual.TempInsanity)
	require.Equal(t, expected.IndefInsanity, actual.IndefInsanity)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}

func sanityRowData(sanity db.SanityState) []any {
	return []any{
		sanity.ID,
		sanity.CharacterID,
		sanity.MaxSanity,
		sanity.CurrentSanity,
		sanity.TempInsanity,
		sanity.IndefInsanity,
		sanity.CreatedAt,
		sanity.UpdatedAt,
	}
}
