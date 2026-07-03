package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetLuck(ctx context.Context, input luckDTO.GetLuckInput) (db.LuckState, error) {
	luck, err := s.next.GetLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) UpsertLuck(ctx context.Context, input luckDTO.UpsertLuckInput) (db.LuckState, error) {
	luck, err := s.next.UpsertLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) DeleteLuck(ctx context.Context, input luckDTO.DeleteLuckInput) error {
	err := s.next.DeleteLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
