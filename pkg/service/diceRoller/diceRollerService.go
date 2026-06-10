package diceRoller

import (
	"context"
	"encoding/json"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
)

type DiceRollerService struct {
	repos *repository.Repository
}

func NewDiceRollerService(repos *repository.Repository) *DiceRollerService {
	return &DiceRollerService{repos: repos}
}

func (s *DiceRollerService) GetLastDiceRolls(ctx context.Context, input db.GetDiceRollsParams) ([]db.DiceRoll, error) {
	diceRolls, err := s.repos.Queries.GetDiceRolls(ctx, input)
	if err != nil {
		return nil, err
	}

	return diceRolls, nil
}

func (s *DiceRollerService) MakeRoll(ctx context.Context, input DiceRollInput) (db.DiceRoll, error) {
	components, err := dice.ParseDiceRollerFormula(input.Formula)
	if err != nil {
		return db.DiceRoll{}, err
	}

	details, result, err := dice.RollDice(components)
	if err != nil {
		return db.DiceRoll{}, err
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return db.DiceRoll{}, err
	}

	diceRoll, err := s.repos.Queries.CreateDiceRoll(ctx, db.CreateDiceRollParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		Expression:  input.Formula,
		Result:      int32(result),
		Details:     detailsJSON,
	})
	if err != nil {
		return db.DiceRoll{}, err
	}

	return diceRoll, nil
}
