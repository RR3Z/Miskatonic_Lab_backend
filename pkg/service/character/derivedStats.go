package character

import (
	"context"
	"errors"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/calculators"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// DerivedStats
func (s *CharacterService) GetDerivedStats(ctx context.Context, input derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.DerivedStat{}, err
	}

	return derivedStats, nil
}

func (s *CharacterService) upsertCalculatedDerivedStats(ctx context.Context, input db.UpsertDerivedStatsParams) error {
	_, err := s.repos.Queries.UpsertDerivedStats(ctx, input)
	return err
}

func (s *CharacterService) clearDerivedStats(ctx context.Context, userID string, characterID pgtype.UUID) error {
	_, err := s.repos.Queries.DeleteDerivedStats(ctx, db.DeleteDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return err
}

func (s *CharacterService) recalculateDerivedStats(
	ctx context.Context,
	userID string,
	characterID pgtype.UUID,
	characteristics db.Characteristic,
	source string,
) {
	if reason, canCalculate := characterHelpers.DerivedStatsRecalculationReadiness(characteristics); !canCalculate {
		if err := s.clearDerivedStats(ctx, userID, characterID); err != nil {
			s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateFailed{
				UserID:      userID,
				CharacterID: characterID.String(),
				Source:      source,
				Err:         err,
			})
			return
		}

		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateSkipped{
			UserID:      userID,
			CharacterID: characterID.String(),
			Source:      source,
			Reason:      reason,
		})
		return
	}

	derivedStatsInput := calculators.CalculateDerivedStats(
		userID,
		characterID,
		characteristics,
	)

	err := s.upsertCalculatedDerivedStats(ctx, derivedStatsInput)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateFailed{
			UserID:      userID,
			CharacterID: characterID.String(),
			Source:      source,
			Err:         err,
		})
		return
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateSucceeded{
		UserID:      userID,
		CharacterID: characterID.String(),
		Source:      source,
	})
}
