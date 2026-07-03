package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetHealth(ctx context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	health, err := s.next.GetHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) UpsertHealth(ctx context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	health, err := s.next.UpsertHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) DeleteHealth(ctx context.Context, input healthDTO.DeleteHealthInput) error {
	err := s.next.DeleteHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
