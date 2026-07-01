package character

import (
	"context"
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
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
	if err := s.validateHealthState(ctx, input); err != nil {
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

func (s *CharacterService) validateHealthState(ctx context.Context, input healthDTO.UpsertHealthInput) error {
	if input.MaxHp != nil && input.CurrentHp != nil {
		if *input.CurrentHp > *input.MaxHp {
			return myErrors.ErrCurrentHealthExceedsMax
		}
		return nil
	}

	if input.MaxHp == nil && input.CurrentHp == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxHp := existing.MaxHp
	if input.MaxHp != nil {
		maxHp = *input.MaxHp
	}

	currentHp := existing.CurrentHp
	if input.CurrentHp != nil {
		currentHp = *input.CurrentHp
	}

	if currentHp > maxHp {
		return myErrors.ErrCurrentHealthExceedsMax
	}

	return nil
}
