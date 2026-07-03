package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetSanity(ctx context.Context, input sanityDTO.GetSanityInput) (db.SanityState, error) {
	sanity, err := s.next.GetSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) UpsertSanity(ctx context.Context, input sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	sanity, err := s.next.UpsertSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) DeleteSanity(ctx context.Context, input sanityDTO.DeleteSanityInput) error {
	err := s.next.DeleteSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
