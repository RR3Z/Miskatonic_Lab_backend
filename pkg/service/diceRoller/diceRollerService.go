package diceRoller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
	"github.com/jackc/pgx/v5"
)

type DiceRollerService struct {
	repos *repository.Repository
}

func NewDiceRollerService(repos *repository.Repository) *DiceRollerService {
	return &DiceRollerService{repos: repos}
}

func (s *DiceRollerService) GetLastDiceRolls(ctx context.Context, input diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error) {
	diceRolls, err := s.repos.Queries.GetDiceRolls(ctx, db.GetDiceRollsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	models := make([]diceRollerDTO.DiceRollModel, len(diceRolls))
	for i, r := range diceRolls {
		models[i] = diceRollerDTO.ToDiceRollModel(r)
	}

	return models, nil
}

func (s *DiceRollerService) MakeRoll(ctx context.Context, input diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error) {
	if err := validateExpression(input.Formula); err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	components, err := dice.ParseDiceRollerFormula(input.Formula)
	if err != nil {
		return diceRollerDTO.DiceRollModel{}, fmt.Errorf("%w: %v", ErrInvalidExpression, err)
	}

	if err := validateComponents(components); err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	details, result, err := dice.RollDice(components)
	if err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	diceRoll, err := s.repos.Queries.CreateDiceRoll(ctx, db.CreateDiceRollParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		Expression:  input.Formula,
		Result:      int32(result),
		Details:     detailsJSON,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return diceRollerDTO.DiceRollModel{}, ErrCharacterNotFound
		}
		return diceRollerDTO.DiceRollModel{}, err
	}

	if err := s.repos.Queries.CleanOldDiceRolls(ctx, db.CleanOldDiceRollsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		slog.Warn("failed to clean old dice rolls", "character_id", input.CharacterID, "error", err)
	}

	return diceRollerDTO.ToDiceRollModel(diceRoll), nil
}
