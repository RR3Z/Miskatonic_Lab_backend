package helpers

import (
	stdErrors "errors"

	appErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/coder/websocket"
)

func UnsupportedCommandTypeEvent(roomID string, actorID string) roomEvents.Event {
	return commandErrorEvent(roomID, actorID, appErrors.AppErrorResponse{
		Code:    appErrors.CodeInvalidRequest,
		Message: "unsupported room command type",
		Details: []appErrors.ErrorDetail{
			appErrors.ValidationDetail("command.type", "unsupported"),
		},
	})
}

func InvalidCommandPayloadEvent(roomID string, actorID string) roomEvents.Event {
	return commandErrorEvent(roomID, actorID, appErrors.AppErrorResponse{
		Code:    appErrors.CodeInvalidRequest,
		Message: "invalid room command payload",
		Details: []appErrors.ErrorDetail{
			appErrors.ParseDetail("command.payload", "invalid_format"),
		},
	})
}

func commandErrorEvent(roomID string, actorID string, payload appErrors.AppErrorResponse) roomEvents.Event {
	return roomEvents.Event{
		Type:    string(roomEvents.EventCommandError),
		RoomID:  roomID,
		ActorID: actorID,
		Payload: payload,
	}
}

func CloseStatusForCommandError(err error) (websocket.StatusCode, string) {
	var closeErr websocket.CloseError
	if stdErrors.As(err, &closeErr) {
		return closeErr.Code, closeErr.Reason
	}

	return websocket.StatusInternalError, "room command failed"
}
