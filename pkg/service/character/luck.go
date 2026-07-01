package character

import (
	"context"
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

// Luck
func (s *CharacterService) GetLuck(ctx context.Context, input luckDTO.GetLuckInput) (db.LuckState, error) {
	luck, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.LuckState{}, err
	}

	return luck, nil
}

func (s *CharacterService) UpsertLuck(ctx context.Context, input luckDTO.UpsertLuckInput) (db.LuckState, error) {
	if err := s.validateLuckState(ctx, input); err != nil {
		return db.LuckState{}, err
	}

	luck, err := s.repos.Queries.UpsertLuckState(ctx, db.UpsertLuckStateParams{
		UserID:       input.UserID,
		CharacterID:  input.CharacterID,
		StartingLuck: input.StartingLuck,
		CurrentLuck:  input.CurrentLuck,
	})
	if err != nil {
		return db.LuckState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return luck, nil
}

func (s *CharacterService) DeleteLuck(ctx context.Context, input luckDTO.DeleteLuckInput) error {
	if _, err := s.repos.Queries.DeleteLuckState(ctx, db.DeleteLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) validateLuckState(ctx context.Context, input luckDTO.UpsertLuckInput) error {
	if input.StartingLuck != nil && input.CurrentLuck != nil {
		if *input.CurrentLuck > *input.StartingLuck {
			return myErrors.ErrCurrentLuckExceedsStarting
		}
		return nil
	}

	if input.StartingLuck == nil && input.CurrentLuck == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	startingLuck := existing.StartingLuck
	if input.StartingLuck != nil {
		startingLuck = *input.StartingLuck
	}

	currentLuck := existing.CurrentLuck
	if input.CurrentLuck != nil {
		currentLuck = *input.CurrentLuck
	}

	if currentLuck > startingLuck {
		return myErrors.ErrCurrentLuckExceedsStarting
	}

	return nil
}
