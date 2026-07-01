package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
)

type FakeEventPublisher struct {
	Events []events.Event
}

func (f *FakeEventPublisher) Publish(_ context.Context, event events.Event) {
	f.Events = append(f.Events, event)
}

type FakeDiceRollerService struct {
	Err   error
	Roll  diceRollerDTO.DiceRollModel
	Rolls []diceRollerDTO.DiceRollModel
}

func (f *FakeDiceRollerService) MakeRoll(_ context.Context, _ diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error) {
	return f.Roll, f.Err
}

func (f *FakeDiceRollerService) GetLastDiceRolls(_ context.Context, _ diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error) {
	return f.Rolls, f.Err
}
