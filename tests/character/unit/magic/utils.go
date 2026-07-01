package tests

import (
	"testing"

	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newTestCharacterServiceForMagic() (*FakeMagicDBTX, *characterServices.CharacterService) {
	dbtx := &FakeMagicDBTX{QueryRowData: magicRowData(testMagicState())}
	repos := &repository.Repository{Queries: db.New(dbtx)}

	return dbtx, characterServices.NewCharacterService(repos)
}

func testMagicState() db.MagicState {
	return db.MagicState{
		ID:          testMagicUUID("22222222-2222-2222-2222-222222222222"),
		CharacterID: testMagicUUID("11111111-1111-1111-1111-111111111111"),
		MaxMp:       10,
		CurrentMp:   7,
		CreatedAt:   testMagicTimestamptz(),
		UpdatedAt:   testMagicTimestamptz(),
	}
}

func testUpsertMagicInput(maxMp *int16, currentMp *int16) magicDTO.UpsertMagicInput {
	return magicDTO.UpsertMagicInput{
		UserID:      "user_1",
		CharacterID: testMagicUUID("11111111-1111-1111-1111-111111111111"),
		MaxMp:       maxMp,
		CurrentMp:   currentMp,
	}
}

func magicInt16(value int16) *int16 {
	return &value
}

func testMagicTimestamptz() pgtype.Timestamptz {
	var value pgtype.Timestamptz
	err := value.Scan("2026-06-07 12:00:00+03")
	if err != nil {
		panic(err)
	}

	return value
}

func testMagicUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func requireSameMagicState(t *testing.T, expected db.MagicState, actual db.MagicState) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxMp, actual.MaxMp)
	require.Equal(t, expected.CurrentMp, actual.CurrentMp)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}

func magicRowData(magic db.MagicState) []any {
	return []any{
		magic.ID,
		magic.CharacterID,
		magic.MaxMp,
		magic.CurrentMp,
		magic.CreatedAt,
		magic.UpdatedAt,
	}
}
