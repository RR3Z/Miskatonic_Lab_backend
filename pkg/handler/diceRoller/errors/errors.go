package diceRollerErrors

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
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

func InvalidInputError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusInternalServerError,
		Message: fallbackMessage,
		Err:     err,
	}
}
