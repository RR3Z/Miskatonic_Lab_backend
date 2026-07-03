package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetDerivedStats(ctx context.Context, input derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.next.GetDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.DerivedStat{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return derivedStats, nil
}

func (s *EventPublishingCharacterService) UpsertDerivedStats(ctx context.Context, input derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.next.UpsertDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.DerivedStat{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return derivedStats, nil
}

func (s *EventPublishingCharacterService) DeleteDerivedStats(ctx context.Context, input derivedStatsDTO.DeleteDerivedStatsInput) error {
	err := s.next.DeleteDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}
