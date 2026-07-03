package helpers

import (
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

func EventFromRoomEventModel(event roomModel.RoomEventModel) roomModel.Event {
	return roomModel.Event{
		Type:    event.Type,
		RoomID:  event.RoomID.String(),
		ActorID: event.ActorID,
		Payload: event.Payload,
	}
}
