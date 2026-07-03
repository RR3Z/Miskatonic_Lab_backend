package tests

import (
	"context"
	"errors"
	"testing"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func diceRollServiceTestUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}

func dbDiceRoll(id, characterID string) db.DiceRoll {
	idUUID := diceRollServiceTestUUID(id)
	charUUID := diceRollServiceTestUUID(characterID)
	return db.DiceRoll{
		ID:          idUUID,
		CharacterID: charUUID,
		UserID:      "user_1",
		Expression:  "1d20",
		Result:      15,
		Details:     []byte(`[{"is_dice":true,"sides":20,"count":1,"result":15}]`),
	}
}

func newServiceWithFakeDBTX(fake *FakeDiceRollerDBTX) *diceRollerServices.DiceRollerService {
	return diceRollerServices.NewDiceRollerService(&repository.Repository{
		Queries: db.New(fake),
	})
}

func TestMakeRoll_CallsCreateDiceRollAndReturnsModel(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")
	roll := dbDiceRoll("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryRowData: []any{
			roll.ID,
			roll.CharacterID,
			roll.UserID,
			roll.Expression,
			roll.Result,
			roll.Details,
			roll.CreatedAt,
		},
	}
	service := newServiceWithFakeDBTX(fake)

	result, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: charID,
		Formula:     "1d20",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fake.QueryRowCalls)
	require.Equal(t, 1, fake.ExecCalls)
	require.Equal(t, roll.ID.String(), result.ID.String())
	require.Equal(t, roll.CharacterID.String(), result.CharacterID.String())
	require.Equal(t, roll.UserID, result.UserID)
	require.Equal(t, roll.Expression, result.Expression)
	require.Equal(t, roll.Result, result.Result)
}

func TestMakeRoll_RoomIDDoesNotChangeCorePersistence(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")
	roomID := diceRollServiceTestUUID("22222222-2222-2222-2222-222222222222")
	roll := dbDiceRoll("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryRowData: []any{
			roll.ID,
			roll.CharacterID,
			roll.UserID,
			roll.Expression,
			roll.Result,
			roll.Details,
			roll.CreatedAt,
		},
	}
	service := newServiceWithFakeDBTX(fake)

	result, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: charID,
		Formula:     "1d20",
		RoomID:      &roomID,
	})

	require.NoError(t, err)
	require.Equal(t, 1, fake.QueryRowCalls)
	require.Equal(t, 1, fake.ExecCalls)
	require.Equal(t, roll.ID.String(), result.ID.String())
}

func TestMakeRoll_CreateDiceRollErrNoRowsMapsToErrCharacterNotFound(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryRowResults: []FakeDiceRollerQueryRowResult{
			{Err: pgx.ErrNoRows},
		},
	}
	service := newServiceWithFakeDBTX(fake)

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: charID,
		Formula:     "1d20",
	})

	require.ErrorIs(t, err, diceRollerServices.ErrCharacterNotFound)
}

func TestMakeRoll_CleanOldDiceRollsErrorDoesNotFailRoll(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")
	roll := dbDiceRoll("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryRowData: []any{
			roll.ID,
			roll.CharacterID,
			roll.UserID,
			roll.Expression,
			roll.Result,
			roll.Details,
			roll.CreatedAt,
		},
		ExecErr: errors.New("cleanup failed"),
	}
	service := newServiceWithFakeDBTX(fake)

	result, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: charID,
		Formula:     "1d20",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fake.ExecCalls)
	require.Equal(t, roll.ID.String(), result.ID.String())
}

func TestGetLastDiceRolls_ReturnsModels(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")
	roll1 := dbDiceRoll("cccccccc-cccc-cccc-cccc-cccccccccccc", "11111111-1111-1111-1111-111111111111")
	roll2 := dbDiceRoll("dddddddd-dddd-dddd-dddd-dddddddddddd", "11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryRows: [][]any{
			{roll1.ID, roll1.CharacterID, roll1.UserID, roll1.Expression, roll1.Result, roll1.Details, roll1.CreatedAt},
			{roll2.ID, roll2.CharacterID, roll2.UserID, roll2.Expression, roll2.Result, roll2.Details, roll2.CreatedAt},
		},
	}
	service := newServiceWithFakeDBTX(fake)

	results, err := service.GetLastDiceRolls(context.Background(), diceRollerDTO.GetLastDiceRollsInput{
		UserID:      "user_1",
		CharacterID: charID,
	})

	require.NoError(t, err)
	require.Equal(t, 1, fake.QueryCalls)
	require.Len(t, results, 2)
	require.Equal(t, roll1.ID.String(), results[0].ID.String())
	require.Equal(t, roll2.ID.String(), results[1].ID.String())
}

func TestGetLastDiceRolls_ReturnsErrorOnQueryFailure(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")

	fake := &FakeDiceRollerDBTX{
		QueryErr: pgx.ErrNoRows,
	}
	service := newServiceWithFakeDBTX(fake)

	_, err := service.GetLastDiceRolls(context.Background(), diceRollerDTO.GetLastDiceRollsInput{
		UserID:      "user_1",
		CharacterID: charID,
	})

	require.Error(t, err)
}
