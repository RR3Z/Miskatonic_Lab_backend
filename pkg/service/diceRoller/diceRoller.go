package diceRoller

import (
	"context"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
)

type IDiceRoller interface {
	MakeRoll(ctx context.Context, input diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error)
	GetLastDiceRolls(ctx context.Context, input diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error)
}
