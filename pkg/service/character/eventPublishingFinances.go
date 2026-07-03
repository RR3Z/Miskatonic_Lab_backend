package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetFinances(ctx context.Context, input financesDTO.GetFinancesInput) (db.Finance, error) {
	finances, err := s.next.GetFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Finance{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return finances, nil
}

func (s *EventPublishingCharacterService) UpsertFinances(ctx context.Context, input financesDTO.UpsertFinancesInput) (db.Finance, error) {
	finances, err := s.next.UpsertFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Finance{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return finances, nil
}

func (s *EventPublishingCharacterService) DeleteFinances(ctx context.Context, input financesDTO.DeleteFinancesInput) error {
	err := s.next.DeleteFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
