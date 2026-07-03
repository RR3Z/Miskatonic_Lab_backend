package tests

import (
	"context"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeDiceRollerHandlerService struct {
	roll      diceRollerDTO.DiceRollModel
	rolls     []diceRollerDTO.DiceRollModel
	err       error
	makeCalls int
	makeInput diceRollerDTO.MakeRollInput
	listCalls int
	listInput diceRollerDTO.GetLastDiceRollsInput
}

func (f *fakeDiceRollerHandlerService) MakeRoll(_ context.Context, input diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error) {
	f.makeCalls++
	f.makeInput = input
	return f.roll, f.err
}

func (f *fakeDiceRollerHandlerService) GetLastDiceRolls(_ context.Context, input diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error) {
	f.listCalls++
	f.listInput = input
	return f.rolls, f.err
}

type fakeRoomAccessChecker struct {
	err error
}

func (f *fakeRoomAccessChecker) EnsureCanPublishRoomEvent(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}
