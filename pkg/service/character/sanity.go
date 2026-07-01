package character

import (
	"context"
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

// Sanity
func (s *CharacterService) GetSanity(ctx context.Context, input sanityDTO.GetSanityInput) (db.SanityState, error) {
	sanity, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.SanityState{}, err
	}

	return sanity, nil
}

func (s *CharacterService) UpsertSanity(ctx context.Context, input sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	if err := s.validateSanityState(ctx, input); err != nil {
		return db.SanityState{}, err
	}

	sanity, err := s.repos.Queries.UpsertSanityState(ctx, db.UpsertSanityStateParams{
		UserID:        input.UserID,
		CharacterID:   input.CharacterID,
		MaxSanity:     input.MaxSanity,
		CurrentSanity: input.CurrentSanity,
		TempInsanity:  input.TempInsanity,
		IndefInsanity: input.IndefInsanity,
	})
	if err != nil {
		return db.SanityState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return sanity, nil
}

func (s *CharacterService) DeleteSanity(ctx context.Context, input sanityDTO.DeleteSanityInput) error {
	if _, err := s.repos.Queries.DeleteSanityState(ctx, db.DeleteSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) validateSanityState(ctx context.Context, input sanityDTO.UpsertSanityInput) error {
	if input.MaxSanity != nil && input.CurrentSanity != nil {
		if *input.CurrentSanity > *input.MaxSanity {
			return myErrors.ErrCurrentSanityExceedsMax
		}
		return nil
	}

	if input.MaxSanity == nil && input.CurrentSanity == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxSanity := existing.MaxSanity
	if input.MaxSanity != nil {
		maxSanity = *input.MaxSanity
	}

	currentSanity := existing.CurrentSanity
	if input.CurrentSanity != nil {
		currentSanity = *input.CurrentSanity
	}

	if currentSanity > maxSanity {
		return myErrors.ErrCurrentSanityExceedsMax
	}

	return nil
}
