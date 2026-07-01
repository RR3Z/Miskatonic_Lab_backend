package diceRoller

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
)

type EventPublishingDiceRollerService struct {
	next      IDiceRoller
	publisher events.EventPublisher
}

func NewEventPublishingDiceRollerService(next IDiceRoller, publisher events.EventPublisher) *EventPublishingDiceRollerService {
	return &EventPublishingDiceRollerService{
		next:      next,
		publisher: publisher,
	}
}

func (s *EventPublishingDiceRollerService) MakeRoll(ctx context.Context, input diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error) {
	roll, err := s.next.MakeRoll(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, diceEvents.DiceRollMakeFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return diceRollerDTO.DiceRollModel{}, err
	}

	s.publisher.Publish(ctx, diceEvents.DiceRollMakeSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		RollID:      roll.ID.String(),
		Expression:  roll.Expression,
		Result:      roll.Result,
		Details:     roll.Details,
	})

	return roll, nil
}

func (s *EventPublishingDiceRollerService) GetLastDiceRolls(ctx context.Context, input diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error) {
	rolls, err := s.next.GetLastDiceRolls(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, diceEvents.DiceRollsListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, diceEvents.DiceRollsListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(rolls),
	})

	return rolls, nil
}
