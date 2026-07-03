package diceRoller

import (
	"context"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
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
	prepared, err := prepareDiceRoll(input.Formula)
	if err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	return s.makeLocalRoll(ctx, input, prepared)
}
