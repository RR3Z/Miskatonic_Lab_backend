package diceRollerErrors

import (
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	serviceDiceRoller "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
)

func InvalidCharacterIDError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "dice.invalid_character_id",
		Message: "invalid character id",
		Err:     err,
	}
}

func InvalidExpressionError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "dice.invalid_expression",
		Message: "invalid dice roll expression",
		Err:     err,
	}
}

func CharacterNotFoundError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusNotFound,
		Code:    "dice.character_not_found",
		Message: "character not found or not owned",
		Err:     err,
	}
}

func InvalidInputError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "dice.invalid_input",
		Message: message,
		Err:     err,
	}
}

func RoomNotAvailableError(err error, message string) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusForbidden,
		Code:    "dice.room_not_available",
		Message: message,
		Err:     err,
	}
}

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, serviceDiceRoller.ErrInvalidExpression):
		return InvalidExpressionError(err)
	case errors.Is(err, serviceDiceRoller.ErrCharacterNotFound):
		return CharacterNotFoundError(err)
	default:
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: fallbackMessage,
			Err:     err,
		}
	}
}
