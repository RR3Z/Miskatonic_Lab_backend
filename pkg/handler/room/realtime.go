package room

import (
	roomDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *RoomHandler) broadcastRoomEvent(eventType roomDTO.EventType, roomID pgtype.UUID, actorID string, payload any) {
	h.hub.Broadcast(roomDTO.Event{
		Type:    string(eventType),
		RoomID:  roomID.String(),
		ActorID: actorID,
		Payload: payload,
	})
}
