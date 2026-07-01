package character

import (
	"context"
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
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
	if err := s.validateMagicState(ctx, input); err != nil {
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

func (s *CharacterService) validateMagicState(ctx context.Context, input magicDTO.UpsertMagicInput) error {
	if input.MaxMp != nil && input.CurrentMp != nil {
		if *input.CurrentMp > *input.MaxMp {
			return myErrors.ErrCurrentMagicExceedsMax
		}
		return nil
	}

	if input.MaxMp == nil && input.CurrentMp == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxMp := existing.MaxMp
	if input.MaxMp != nil {
		maxMp = *input.MaxMp
	}

	currentMp := existing.CurrentMp
	if input.CurrentMp != nil {
		currentMp = *input.CurrentMp
	}

	if currentMp > maxMp {
		return myErrors.ErrCurrentMagicExceedsMax
	}

	return nil
}
