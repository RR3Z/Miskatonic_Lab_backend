package diceRoller

import (
	"fmt"
	"strings"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
)

func validateExpression(expression string) error {
	cleaned := strings.TrimSpace(expression)
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if cleaned == "" {
		return fmt.Errorf("%w: expression is required", ErrInvalidExpression)
	}

	if strings.Contains(cleaned, "-") {
		return fmt.Errorf("%w: subtraction is not allowed", ErrInvalidExpression)
	}

	if !strings.Contains(cleaned, "d") {
		return fmt.Errorf("%w: expression must contain at least one dice term", ErrInvalidExpression)
	}

	return nil
}

func validateComponents(components []dice.DiceRollFormulaComponent) error {
	hasDice := false
	for _, c := range components {
		if c.IsDice {
			hasDice = true
			if c.Count <= 0 {
				return fmt.Errorf("%w: dice count must be greater than 0", ErrInvalidExpression)
			}
			if c.Sides <= 0 {
				return fmt.Errorf("%w: dice sides must be greater than 0", ErrInvalidExpression)
			}
		} else {
			if c.Count < 0 {
				return fmt.Errorf("%w: negative modifiers are not allowed", ErrInvalidExpression)
			}
		}
	}
	if !hasDice {
		return fmt.Errorf("%w: expression must contain at least one dice term", ErrInvalidExpression)
	}
	return nil
}
