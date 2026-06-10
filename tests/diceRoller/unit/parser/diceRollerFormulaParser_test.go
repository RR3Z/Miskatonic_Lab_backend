package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
	"github.com/stretchr/testify/require"
)

func TestParseDiceRollerFormula_SimpleDice(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)
}

func TestParseDiceRollerFormula_SingleDieWithoutCount(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("d20")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 20, components[0].Sides)
	require.Equal(t, 0, components[0].Count)
}

func TestParseDiceRollerFormula_NegativeDie(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("-d6")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, -1, components[0].Count)
}

func TestParseDiceRollerFormula_NegativeMultipleDice(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("-3d6")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, -3, components[0].Count)
}

func TestParseDiceRollerFormula_PlainModifier(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("5")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.False(t, components[0].IsDice)
	require.Equal(t, 5, components[0].Count)
}

func TestParseDiceRollerFormula_NegativeModifier(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("-5")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.False(t, components[0].IsDice)
	require.Equal(t, -5, components[0].Count)
}

func TestParseDiceRollerFormula_DicePlusModifier(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6+2")
	require.NoError(t, err)
	require.Len(t, components, 2)

	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)

	require.False(t, components[1].IsDice)
	require.Equal(t, 2, components[1].Count)
}

func TestParseDiceRollerFormula_DiceMinusModifier(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6-2")
	require.NoError(t, err)
	require.Len(t, components, 2)

	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)

	require.False(t, components[1].IsDice)
	require.Equal(t, -2, components[1].Count)
}

func TestParseDiceRollerFormula_MultipleDiceTypes(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("2d8+1d4")
	require.NoError(t, err)
	require.Len(t, components, 2)

	require.True(t, components[0].IsDice)
	require.Equal(t, 8, components[0].Sides)
	require.Equal(t, 2, components[0].Count)

	require.True(t, components[1].IsDice)
	require.Equal(t, 4, components[1].Sides)
	require.Equal(t, 1, components[1].Count)
}

func TestParseDiceRollerFormula_MultipleModifiers(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("1d20+5+3")
	require.NoError(t, err)
	require.Len(t, components, 3)

	require.True(t, components[0].IsDice)
	require.Equal(t, 20, components[0].Sides)
	require.Equal(t, 1, components[0].Count)

	require.False(t, components[1].IsDice)
	require.Equal(t, 5, components[1].Count)

	require.False(t, components[2].IsDice)
	require.Equal(t, 3, components[2].Count)
}

func TestParseDiceRollerFormula_MixedDiceAndModifiers(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6-2d4+5")
	require.NoError(t, err)
	require.Len(t, components, 3)

	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)

	require.True(t, components[1].IsDice)
	require.Equal(t, 4, components[1].Sides)
	require.Equal(t, -2, components[1].Count)

	require.False(t, components[2].IsDice)
	require.Equal(t, 5, components[2].Count)
}

func TestParseDiceRollerFormula_HandlesSpaces(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6 + 2")
	require.NoError(t, err)
	require.Len(t, components, 2)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)
	require.False(t, components[1].IsDice)
	require.Equal(t, 2, components[1].Count)
}

func TestParseDiceRollerFormula_AlreadyNormalizedDoubleNegative(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("3d6+-2")
	require.NoError(t, err)
	require.Len(t, components, 2)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)
	require.False(t, components[1].IsDice)
	require.Equal(t, -2, components[1].Count)
}

func TestParseDiceRollerFormula_EmptyString(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("")
	require.NoError(t, err)
	require.Empty(t, components)
}

func TestParseDiceRollerFormula_WhitespaceOnly(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("   ")
	require.NoError(t, err)
	require.Empty(t, components)
}

func TestParseDiceRollerFormula_ErrorInvalidDiceCount(t *testing.T) {
	_, err := dice.ParseDiceRollerFormula("xd6")
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong amount of dices")
}

func TestParseDiceRollerFormula_ErrorInvalidDiceSides(t *testing.T) {
	_, err := dice.ParseDiceRollerFormula("3dx")
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong amount of dice sides")
}

func TestParseDiceRollerFormula_ErrorInvalidModifier(t *testing.T) {
	_, err := dice.ParseDiceRollerFormula("3d6+abc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong modifier value")
}

func TestParseDiceRollerFormula_ErrorPercentileNotSupported(t *testing.T) {
	_, err := dice.ParseDiceRollerFormula("d%")
	require.Error(t, err)
}

func TestParseDiceRollerFormula_ZeroCountDie(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("0d6")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 0, components[0].Count)
}

func TestParseDiceRollerFormula_SingleDieWithLeadingPlus(t *testing.T) {
	components, err := dice.ParseDiceRollerFormula("+3d6")
	require.NoError(t, err)
	require.Len(t, components, 1)
	require.True(t, components[0].IsDice)
	require.Equal(t, 6, components[0].Sides)
	require.Equal(t, 3, components[0].Count)
}


