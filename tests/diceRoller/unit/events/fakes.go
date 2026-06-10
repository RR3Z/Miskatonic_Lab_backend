package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
)

type FakeEventPublisher struct {
	Events []events.Event
}

func (f *FakeEventPublisher) Publish(_ context.Context, event events.Event) {
	f.Events = append(f.Events, event)
}

type FakeDiceRollerService struct {
	Err   error
	Roll  db.DiceRoll
	Rolls []db.DiceRoll
}

func (f *FakeDiceRollerService) MakeRoll(ctx context.Context, input diceRoller.DiceRollInput) (db.DiceRoll, error) {
	return f.Roll, f.Err
}

func (f *FakeDiceRollerService) GetLastDiceRolls(ctx context.Context, input db.GetDiceRollsParams) ([]db.DiceRoll, error) {
	return f.Rolls, f.Err
}
