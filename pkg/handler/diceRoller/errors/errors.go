package diceRollerErrors

import (
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	serviceDiceRoller "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
)

func InvalidCharacterIDError(err error) *myErrors.AppError {
	return myErrors.NewAppError("dice.invalid_character_id", err)
}

func InvalidExpressionError(err error) *myErrors.AppError {
	return myErrors.NewAppError("dice.invalid_expression", err)
}

func CharacterNotFoundError(err error) *myErrors.AppError {
	return myErrors.NewAppError("dice.character_not_found", err)
}

func InvalidInputError(_ string, err error) *myErrors.AppError {
	return myErrors.NewAppError("dice.invalid_input", err)
}

func RoomNotAvailableError(err error, _ string) *myErrors.AppError {
	return myErrors.NewAppError("dice.room_not_available", err)
}

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, serviceDiceRoller.ErrInvalidExpression):
		return InvalidExpressionError(err)
	case errors.Is(err, serviceDiceRoller.ErrCharacterNotFound):
		return CharacterNotFoundError(err)
	default:
		_ = fallbackMessage
		return myErrors.NewAppError(myErrors.CodeInternalError, err)
	}
}
