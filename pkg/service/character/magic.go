package character

import (
	"context"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

// Magic
func (s *CharacterService) GetMagic(ctx context.Context, input magicDTO.GetMagicInput) (db.MagicState, error) {
	magic, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.MagicState{}, err
	}

	return magic, nil
}

func (s *CharacterService) UpsertMagic(ctx context.Context, input magicDTO.UpsertMagicInput) (db.MagicState, error) {
	if err := validateNonNegative(characterErrors.ErrStateNegative, input.MaxMp, input.CurrentMp); err != nil {
		return db.MagicState{}, err
	}

	if err := s.validateStateMax(ctx, input.MaxMp, input.CurrentMp,
		func(ctx context.Context) (int16, int16, error) {
			existing, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
				UserID: input.UserID, CharacterID: input.CharacterID,
			})
			if err != nil {
				return 0, 0, err
			}
			return existing.MaxMp, existing.CurrentMp, nil
		},
		myErrors.ErrCurrentMagicExceedsMax,
	); err != nil {
		return db.MagicState{}, err
	}

	magic, err := s.repos.Queries.UpsertMagicState(ctx, db.UpsertMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		MaxMp:       input.MaxMp,
		CurrentMp:   input.CurrentMp,
	})
	if err != nil {
		return db.MagicState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return magic, nil
}

func (s *CharacterService) DeleteMagic(ctx context.Context, input magicDTO.DeleteMagicInput) error {
	if _, err := s.repos.Queries.DeleteMagicState(ctx, db.DeleteMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}
