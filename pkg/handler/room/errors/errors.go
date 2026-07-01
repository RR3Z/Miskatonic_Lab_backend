package roomErrors

import (
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	serviceroom "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
)

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, serviceroom.ErrInvalidInput):
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Code:    "room.invalid_input",
			Message: "invalid room input",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrRoomNotFound):
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Code:    "room.not_found",
			Message: "room not found",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrNotMember):
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Code:    "room.not_member",
			Message: "not a member of this room",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrNotOwner):
		return &myErrors.AppError{
			Status:  http.StatusForbidden,
			Code:    "room.not_owner",
			Message: "only the room owner can perform this action",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrRoomFull):
		return &myErrors.AppError{
			Status:  http.StatusConflict,
			Code:    "room.full",
			Message: "room is full",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrAlreadyMember):
		return &myErrors.AppError{
			Status:  http.StatusConflict,
			Code:    "room.already_member",
			Message: "already a member of this room",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrCannotKickOwner):
		return &myErrors.AppError{
			Status:  http.StatusForbidden,
			Code:    "room.cannot_kick_owner",
			Message: "cannot kick the room owner",
			Err:     err,
		}
	case errors.Is(err, serviceroom.ErrCharacterNotOwned):
		return &myErrors.AppError{
			Status:  http.StatusForbidden,
			Code:    "room.character_not_owned",
			Message: "character does not belong to you",
			Err:     err,
		}
	default:
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: fallbackMessage,
			Err:     err,
		}
	}
}

func InvalidIDError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "room.invalid_id",
		Message: "invalid room id",
		Err:     err,
	}
}

func InvalidInputError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "room.invalid_input",
		Message: message,
		Err:     err,
	}
}
