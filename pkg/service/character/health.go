package character

import (
	"context"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

// Health
func (s *CharacterService) GetHealth(ctx context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	health, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.HealthState{}, err
	}

	return health, nil
}

func (s *CharacterService) UpsertHealth(ctx context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	if err := validateNonNegative(characterErrors.ErrStateNegative, input.MaxHp, input.CurrentHp); err != nil {
		return db.HealthState{}, err
	}

	if err := s.validateStateMax(ctx, input.MaxHp, input.CurrentHp,
		func(ctx context.Context) (int16, int16, error) {
			existing, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
				UserID: input.UserID, CharacterID: input.CharacterID,
			})
			if err != nil {
				return 0, 0, err
			}
			return existing.MaxHp, existing.CurrentHp, nil
		},
		myErrors.ErrCurrentHealthExceedsMax,
	); err != nil {
		return db.HealthState{}, err
	}

	health, err := s.repos.Queries.UpsertHealthState(ctx, db.UpsertHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		MaxHp:       input.MaxHp,
		CurrentHp:   input.CurrentHp,
		MajorWound:  input.MajorWound,
		Unconscious: input.Unconscious,
		Dying:       input.Dying,
		Dead:        input.Dead,
	})
	if err != nil {
		return db.HealthState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return health, nil
}

func (s *CharacterService) DeleteHealth(ctx context.Context, input healthDTO.DeleteHealthInput) error {
	if _, err := s.repos.Queries.DeleteHealthState(ctx, db.DeleteHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}
