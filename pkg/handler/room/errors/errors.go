package roomErrors

import (
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	serviceroom "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
)

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, serviceroom.ErrInvalidInput), errors.Is(err, serviceroom.ErrInvalidPassword):
		return myErrors.NewAppError("room.invalid_input", err)
	case errors.Is(err, serviceroom.ErrRoomNotFound):
		return myErrors.NewAppError("room.not_found", err)
	case errors.Is(err, serviceroom.ErrNotMember):
		return myErrors.NewAppError("room.not_member", err)
	case errors.Is(err, serviceroom.ErrNotOwner):
		return myErrors.NewAppError("room.not_owner", err)
	case errors.Is(err, serviceroom.ErrRoomFull):
		return myErrors.NewAppError("room.full", err)
	case errors.Is(err, serviceroom.ErrAlreadyMember):
		return myErrors.NewAppError("room.already_member", err)
	case errors.Is(err, serviceroom.ErrCannotKickOwner):
		return myErrors.NewAppError("room.cannot_kick_owner", err)
	case errors.Is(err, serviceroom.ErrCharacterNotOwned):
		return myErrors.NewAppError("room.character_not_owned", err)
	default:
		_ = fallbackMessage
		return myErrors.NewAppError(myErrors.CodeInternalError, err)
	}
}

func InvalidIDError(err error) *myErrors.AppError {
	return myErrors.NewAppError("room.invalid_id", err)
}

func InvalidInputError(_ string, err error) *myErrors.AppError {
	return myErrors.NewAppError("room.invalid_input", err)
}
