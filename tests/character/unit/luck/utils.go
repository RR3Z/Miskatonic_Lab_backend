package tests

import (
	"testing"

	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newTestCharacterServiceForLuck() (*FakeLuckDBTX, *characterServices.CharacterService) {
	dbtx := &FakeLuckDBTX{QueryRowData: luckRowData(testLuckState())}
	repos := &repository.Repository{Queries: db.New(dbtx)}

	return dbtx, characterServices.NewCharacterService(repos)
}

func testLuckState() db.LuckState {
	return db.LuckState{
		ID:           testLuckUUID("22222222-2222-2222-2222-222222222222"),
		CharacterID:  testLuckUUID("11111111-1111-1111-1111-111111111111"),
		StartingLuck: 60,
		CurrentLuck:  40,
		CreatedAt:    testLuckTimestamptz(),
		UpdatedAt:    testLuckTimestamptz(),
	}
}

func testUpsertLuckInput(startingLuck *int16, currentLuck *int16) luckDTO.UpsertLuckInput {
	return luckDTO.UpsertLuckInput{
		UserID:       "user_1",
		CharacterID:  testLuckUUID("11111111-1111-1111-1111-111111111111"),
		StartingLuck: startingLuck,
		CurrentLuck:  currentLuck,
	}
}

func luckInt16(value int16) *int16 {
	return &value
}

func testLuckTimestamptz() pgtype.Timestamptz {
	var value pgtype.Timestamptz
	err := value.Scan("2026-06-07 12:00:00+03")
	if err != nil {
		panic(err)
	}

	return value
}

func testLuckUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func requireSameLuckState(t *testing.T, expected db.LuckState, actual db.LuckState) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.StartingLuck, actual.StartingLuck)
	require.Equal(t, expected.CurrentLuck, actual.CurrentLuck)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}

func luckRowData(luck db.LuckState) []any {
	return []any{
		luck.ID,
		luck.CharacterID,
		luck.StartingLuck,
		luck.CurrentLuck,
		luck.CreatedAt,
		luck.UpdatedAt,
	}
}
