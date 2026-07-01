package diceRoller

import "fmt"

func validateExpressionNotEmpty(expression string) error {
	if expression == "" {
		return fmt.Errorf("%w: expression is required", ErrInvalidExpression)
	}
	return nil
}
