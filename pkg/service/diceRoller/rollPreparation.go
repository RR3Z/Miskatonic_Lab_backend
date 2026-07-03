package diceRoller

import (
	"encoding/json"
	"fmt"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
)

type preparedDiceRoll struct {
	detailsJSON []byte
	result      int
}

func prepareDiceRoll(formula string) (preparedDiceRoll, error) {
	if err := validateExpression(formula); err != nil {
		return preparedDiceRoll{}, err
	}

	components, err := dice.ParseDiceRollerFormula(formula)
	if err != nil {
		return preparedDiceRoll{}, fmt.Errorf("%w: %v", ErrInvalidExpression, err)
	}

	if err := validateComponents(components); err != nil {
		return preparedDiceRoll{}, err
	}

	details, result, err := dice.RollDice(components)
	if err != nil {
		return preparedDiceRoll{}, err
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return preparedDiceRoll{}, err
	}

	return preparedDiceRoll{detailsJSON: detailsJSON, result: result}, nil
}
