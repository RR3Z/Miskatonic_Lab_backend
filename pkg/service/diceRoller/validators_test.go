package diceRoller

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
)

func TestValidateExpression_ValidFormulas(t *testing.T) {
	valid := []string{
		"1d20",
		"1d20+5",
		"2d6+1d4",
		"2d6+1d4+5",
		"3d6 + 2",
	}
	for _, formula := range valid {
		t.Run(formula, func(t *testing.T) {
			if err := validateExpression(formula); err != nil {
				t.Errorf("expected valid, got: %v", err)
			}
		})
	}
}

func TestValidateExpression_InvalidFormulas(t *testing.T) {
	invalid := []string{
		"",
		"   ",
		"5",
		"5+3",
		"3d6-1d4",
		"1d20-5",
		"3d6+-2",
	}
	for _, formula := range invalid {
		t.Run(formula, func(t *testing.T) {
			if err := validateExpression(formula); err == nil {
				t.Errorf("expected error for %q", formula)
			}
		})
	}
}

func TestValidateExpression_Empty(t *testing.T) {
	if err := validateExpression(""); err == nil {
		t.Fatal("expected error for empty")
	}
}

func TestValidateExpression_WhitespaceOnly(t *testing.T) {
	if err := validateExpression("   "); err == nil {
		t.Fatal("expected error for whitespace only")
	}
}

func TestValidateComponents_ValidComponents(t *testing.T) {
	components := []dice.DiceRollFormulaComponent{
		{IsDice: true, Count: 2, Sides: 6},
		{IsDice: true, Count: 1, Sides: 4},
		{IsDice: false, Count: 5},
	}
	if err := validateComponents(components); err != nil {
		t.Errorf("expected valid, got: %v", err)
	}
}

func TestValidateComponents_RejectsZeroDiceCount(t *testing.T) {
	components := []dice.DiceRollFormulaComponent{
		{IsDice: true, Count: 0, Sides: 6},
	}
	if err := validateComponents(components); err == nil {
		t.Fatal("expected error for zero count")
	}
}

func TestValidateComponents_RejectsZeroDiceSides(t *testing.T) {
	components := []dice.DiceRollFormulaComponent{
		{IsDice: true, Count: 1, Sides: 0},
	}
	if err := validateComponents(components); err == nil {
		t.Fatal("expected error for zero sides")
	}
}

func TestValidateComponents_RejectsNegativeModifier(t *testing.T) {
	components := []dice.DiceRollFormulaComponent{
		{IsDice: true, Count: 1, Sides: 20},
		{IsDice: false, Count: -3},
	}
	if err := validateComponents(components); err == nil {
		t.Fatal("expected error for negative modifier")
	}
}

func TestValidateComponents_RequiresAtLeastOneDice(t *testing.T) {
	components := []dice.DiceRollFormulaComponent{
		{IsDice: false, Count: 5},
		{IsDice: false, Count: 3},
	}
	if err := validateComponents(components); err == nil {
		t.Fatal("expected error for no dice terms")
	}
}
