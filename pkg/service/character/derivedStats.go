package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/calculators"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/helpers"
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

func (s *CharacterService) UpsertDerivedStats(ctx context.Context, input derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	if err := validateNonNegative(characterErrors.ErrDerivedStatsNegative, input.Speed, input.DodgeValue); err != nil {
		return db.DerivedStat{}, err
	}
	if input.DamageBonus != nil {
		bonus := *input.DamageBonus
		valid := false
		if bonus == "-2" || bonus == "-1" || bonus == "0" {
			valid = true
		} else if len(bonus) > 1 && bonus[0] == '+' {
			valid = true
		}
		if !valid {
			return db.DerivedStat{}, characterErrors.ErrInvalidDamageBonus
		}
	}

	derivedStats, err := s.repos.Queries.UpsertDerivedStats(ctx, db.UpsertDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		Speed:       input.Speed,
		Physique:    input.Physique,
		DamageBonus: input.DamageBonus,
		DodgeValue:  input.DodgeValue,
	})
	if err != nil {
		return db.DerivedStat{}, characterErrors.MapCharacterConstraintError(err)
	}

	return derivedStats, nil
}

func (s *CharacterService) DeleteDerivedStats(ctx context.Context, input derivedStatsDTO.DeleteDerivedStatsInput) error {
	if _, err := s.repos.Queries.DeleteDerivedStats(ctx, db.DeleteDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) getCharacteristicsForDerivedStatsRecalculation(
	ctx context.Context,
	userID string,
	characterID pgtype.UUID,
) (db.Characteristic, bool) {
	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return db.Characteristic{}, false
	}

	return characteristics, true
}

func (s *CharacterService) recalculateDerivedStats(
	ctx context.Context,
	userID string,
	characterID pgtype.UUID,
	age *int16,
	characteristics db.Characteristic,
	source string,
) {
	if reason, canCalculate := characterHelpers.DerivedStatsRecalculationReadiness(age, characteristics); !canCalculate {
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
		*age,
		characteristics,
	)

	_, err := s.UpsertDerivedStats(ctx, derivedStatsDTO.UpsertDerivedStatsInput{
		UserID:      derivedStatsInput.UserID,
		CharacterID: derivedStatsInput.CharacterID,
		Speed:       derivedStatsInput.Speed,
		Physique:    derivedStatsInput.Physique,
		DamageBonus: derivedStatsInput.DamageBonus,
		DodgeValue:  derivedStatsInput.DodgeValue,
	})
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
