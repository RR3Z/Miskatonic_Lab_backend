package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetMagic(ctx context.Context, input magicDTO.GetMagicInput) (db.MagicState, error) {
	magic, err := s.next.GetMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) UpsertMagic(ctx context.Context, input magicDTO.UpsertMagicInput) (db.MagicState, error) {
	magic, err := s.next.UpsertMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) DeleteMagic(ctx context.Context, input magicDTO.DeleteMagicInput) error {
	err := s.next.DeleteMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
