package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type DiceRollFormulaComponent struct {
	IsDice bool
	Sides  int
	Count  int
}

func ParseDiceRollerFormula(formula string) ([]DiceRollFormulaComponent, error) {
	formula = prepareFormulaForParsing(formula)

	parts := strings.Split(formula, "+")

	var result []DiceRollFormulaComponent
	for _, part := range parts {
		var count int
		var err error

		if part == "" {
			continue
		}

		left, right, found := strings.Cut(part, "d")
		if found {
			if left != "" {
				if left == "-" {
					count = -1
				} else {
					count, err = strconv.Atoi(left)
					if err != nil {
						return nil, fmt.Errorf("wrong amount of dices: %s", left)
					}
				}

			}

			sides, err := strconv.Atoi(right)
			if err != nil {
				return nil, fmt.Errorf("wrong amount of dice sides: %s", right)
			}

			result = append(result, DiceRollFormulaComponent{
				IsDice: true,
				Sides:  sides,
				Count:  count,
			})
		} else {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("wrong modifier value: %s", part)
			}

			result = append(result, DiceRollFormulaComponent{IsDice: false, Count: val})
		}
	}

	return result, nil
}

func prepareFormulaForParsing(formula string) string {
	// Remove spaces
	formula = strings.TrimSpace(formula)
	formula = strings.ReplaceAll(formula, " ", "")

	// Negative Modifiers
	formula = strings.ReplaceAll(formula, "+-", "-")
	formula = strings.ReplaceAll(formula, "-", "+-")

	return formula
}
