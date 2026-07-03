package helpers

import (
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

func EventFromRoomEventModel(event roomModel.RoomEventModel) roomEvents.Event {
	return roomEvents.Event{
		Type:    event.Type,
		RoomID:  event.RoomID.String(),
		ActorID: event.ActorID,
		Payload: event.Payload,
	}
}
