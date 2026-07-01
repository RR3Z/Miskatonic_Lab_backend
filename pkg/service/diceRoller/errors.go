package diceRoller

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid dice roll expression")
	ErrCharacterNotFound = errors.New("character not found or not owned by user")
)
