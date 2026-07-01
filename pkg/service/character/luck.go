package character

import (
	"context"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
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
	if err := validateNonNegative(characterErrors.ErrStateNegative, input.StartingLuck, input.CurrentLuck); err != nil {
		return db.LuckState{}, err
	}

	if err := s.validateStateMax(ctx, input.StartingLuck, input.CurrentLuck,
		func(ctx context.Context) (int16, int16, error) {
			existing, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
				UserID: input.UserID, CharacterID: input.CharacterID,
			})
			if err != nil {
				return 0, 0, err
			}
			return existing.StartingLuck, existing.CurrentLuck, nil
		},
		myErrors.ErrCurrentLuckExceedsStarting,
	); err != nil {
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
