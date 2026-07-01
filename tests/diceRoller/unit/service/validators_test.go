package tests

import (
	"context"
	"testing"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/stretchr/testify/require"
)

func TestMakeRoll_ValidFormulaReachesPersistence(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")

	cases := []string{
		"1d20",
		"1d20+5",
		"2d6+1d4",
		"2d6+1d4+5",
		"3d6 + 2",
	}

	for _, formula := range cases {
		t.Run(formula, func(t *testing.T) {
			roll := dbDiceRoll("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "11111111-1111-1111-1111-111111111111")
			roll.Expression = formula
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
				Formula:     formula,
			})

			require.NoError(t, err)
			require.Equal(t, 1, fake.QueryRowCalls)
			require.Equal(t, 1, fake.ExecCalls)
		})
	}
}

func TestMakeRoll_InvalidFormulaFailsBeforePersistence(t *testing.T) {
	charID := diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111")

	cases := []string{
		"",
		"   ",
		"5",
		"5+3",
		"3d6-1d4",
		"1d20-5",
		"3d6+-2",
		"0d6",
		"1d0",
		"d20",
		"xd6",
	}

	for _, formula := range cases {
		t.Run(formula, func(t *testing.T) {
			fake := &FakeDiceRollerDBTX{}
			service := newServiceWithFakeDBTX(fake)

			_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
				UserID:      "user_1",
				CharacterID: charID,
				Formula:     formula,
			})

			require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
			require.Zero(t, fake.QueryRowCalls)
			require.Zero(t, fake.QueryCalls)
			require.Zero(t, fake.ExecCalls)
		})
	}
}
