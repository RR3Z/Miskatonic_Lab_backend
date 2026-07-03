package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetCharacteristics(ctx context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	characteristics, err := s.next.GetCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) UpsertCharacteristics(ctx context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	characteristics, err := s.next.UpsertCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) DeleteCharacteristics(ctx context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error {
	err := s.next.DeleteCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
