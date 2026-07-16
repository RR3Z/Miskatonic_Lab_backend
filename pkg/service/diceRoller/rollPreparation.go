package diceRoller

import (
	"encoding/json"
	"fmt"
	"strings"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
)

type preparedDiceRoll struct {
	detailsJSON []byte
	result      int
}

type formulaRollDetails struct {
	Rolls []dice.DiceRollDetail `json:"rolls"`
}

type d100RollDetails struct {
	Mode       diceRollerDTO.D100Mode `json:"mode"`
	Units      int                    `json:"units"`
	Tens       []int                  `json:"tens"`
	Candidates []int                  `json:"candidates"`
	Selected   int                    `json:"selected"`
}

func prepareDiceRoll(formula string, d100Mode *diceRollerDTO.D100Mode) (preparedDiceRoll, error) {
	if d100Mode != nil {
		return prepareD100Roll(formula, *d100Mode)
	}

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

	detailsJSON, err := json.Marshal(formulaRollDetails{Rolls: details})
	if err != nil {
		return preparedDiceRoll{}, err
	}

	return preparedDiceRoll{detailsJSON: detailsJSON, result: result}, nil
}

func prepareD100Roll(formula string, mode diceRollerDTO.D100Mode) (preparedDiceRoll, error) {
	if !mode.IsValid() {
		return preparedDiceRoll{}, fmt.Errorf("%w: d100 mode is invalid", ErrInvalidExpression)
	}
	if normalizeDiceFormula(formula) != "1d100" {
		return preparedDiceRoll{}, fmt.Errorf("%w: d100 mode requires expression 1d100", ErrInvalidExpression)
	}

	tensCount := 1
	if mode != diceRollerDTO.D100ModeNormal {
		tensCount = 2
	}

	roll, err := dice.RollD100(tensCount)
	if err != nil {
		return preparedDiceRoll{}, err
	}

	selected := roll.Candidates[0]
	for _, candidate := range roll.Candidates[1:] {
		if mode == diceRollerDTO.D100ModeBonus && candidate < selected {
			selected = candidate
		}
		if mode == diceRollerDTO.D100ModePenalty && candidate > selected {
			selected = candidate
		}
	}

	detailsJSON, err := json.Marshal(d100RollDetails{
		Mode:       mode,
		Units:      roll.Units,
		Tens:       roll.Tens,
		Candidates: roll.Candidates,
		Selected:   selected,
	})
	if err != nil {
		return preparedDiceRoll{}, err
	}

	return preparedDiceRoll{detailsJSON: detailsJSON, result: selected}, nil
}

func normalizeDiceFormula(formula string) string {
	return strings.ReplaceAll(strings.TrimSpace(formula), " ", "")
}
