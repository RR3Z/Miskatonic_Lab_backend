package room

import (
	roomDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
)

func (h *RoomHandler) broadcastRoomEvents(events []roomDTO.RoomEventModel) {
	for _, event := range events {
		h.hub.Broadcast(wsHelpers.EventFromRoomEventModel(event))
	}
}
