package tests

import (
	"context"
	"encoding/json"
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
		Details:     []byte(`{"rolls":[{"type":"dice","sides":20,"rolls":[15]}]}`),
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
	require.JSONEq(t, `{"rolls":[{"type":"dice","sides":20,"rolls":[15]}]}`, string(result.Details))

	detailsJSON, ok := fake.LastQueryRowArgs[2].([]byte)
	require.True(t, ok)
	var details struct {
		Rolls []map[string]any `json:"rolls"`
	}
	require.NoError(t, json.Unmarshal(detailsJSON, &details))
	require.Len(t, details.Rolls, 1)
	require.Equal(t, "dice", details.Rolls[0]["type"])
	require.Equal(t, float64(20), details.Rolls[0]["sides"])
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

func TestMakeRoll_D100ModesStoreStructuredDetails(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")

	for _, mode := range []diceRollerDTO.D100Mode{
		diceRollerDTO.D100ModeNormal,
		diceRollerDTO.D100ModeBonus,
		diceRollerDTO.D100ModePenalty,
	} {
		t.Run(string(mode), func(t *testing.T) {
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

			_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
				UserID:      "user_1",
				CharacterID: charID,
				Formula:     "1d100",
				D100Mode:    &mode,
			})
			require.NoError(t, err)
			require.Len(t, fake.LastQueryRowArgs, 5)

			detailsJSON, ok := fake.LastQueryRowArgs[2].([]byte)
			require.True(t, ok)

			var details struct {
				Mode       diceRollerDTO.D100Mode `json:"mode"`
				Tens       []int                  `json:"tens"`
				Candidates []int                  `json:"candidates"`
				Selected   int                    `json:"selected"`
			}
			require.NoError(t, json.Unmarshal(detailsJSON, &details))
			require.Equal(t, mode, details.Mode)
			require.Len(t, details.Tens, map[diceRollerDTO.D100Mode]int{
				diceRollerDTO.D100ModeNormal:  1,
				diceRollerDTO.D100ModeBonus:   2,
				diceRollerDTO.D100ModePenalty: 2,
			}[mode])
			require.Equal(t, details.Selected, int(fake.LastQueryRowArgs[1].(int32)))
			require.Contains(t, details.Candidates, details.Selected)
		})
	}
}

func TestMakeRoll_D100ModeRejectsOtherExpressionsAndUnknownModes(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")
	unknown := diceRollerDTO.D100Mode("advantage")

	for _, input := range []diceRollerDTO.MakeRollInput{
		{UserID: "user_1", CharacterID: charID, Formula: "1d6", D100Mode: d100ModePtr(diceRollerDTO.D100ModeBonus)},
		{UserID: "user_1", CharacterID: charID, Formula: "1d100", D100Mode: &unknown},
	} {
		fake := &FakeDiceRollerDBTX{}
		service := newServiceWithFakeDBTX(fake)

		_, err := service.MakeRoll(context.Background(), input)

		require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
		require.Zero(t, fake.QueryRowCalls)
	}
}

func TestToDiceRollModelExposesPersistedStructuredDetails(t *testing.T) {
	roll := dbDiceRoll("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "11111111-1111-1111-1111-111111111111")
	roll.Expression = "1d100"
	roll.Result = 24
	roll.Details = []byte(`{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`)

	model := diceRollerDTO.ToDiceRollModel(roll)

	require.JSONEq(t, `{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`, string(model.Details))
}

func d100ModePtr(mode diceRollerDTO.D100Mode) *diceRollerDTO.D100Mode {
	return &mode
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
