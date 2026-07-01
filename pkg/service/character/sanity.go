package character

import (
	"context"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
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
	if err := validateNonNegative(characterErrors.ErrStateNegative, input.MaxSanity, input.CurrentSanity); err != nil {
		return db.SanityState{}, err
	}

	if err := s.validateStateMax(ctx, input.MaxSanity, input.CurrentSanity,
		func(ctx context.Context) (int16, int16, error) {
			existing, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
				UserID: input.UserID, CharacterID: input.CharacterID,
			})
			if err != nil {
				return 0, 0, err
			}
			return existing.MaxSanity, existing.CurrentSanity, nil
		},
		myErrors.ErrCurrentSanityExceedsMax,
	); err != nil {
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
