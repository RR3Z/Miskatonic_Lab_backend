package room

import (
	roomDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *RoomHandler) CloseDeletedRooms(roomIDs []pgtype.UUID, reason string) {
	for _, roomID := range roomIDs {
		h.presence.ForgetRoom(roomID)
		h.closeRoom(roomID, reason)
	}
}

func (h *RoomHandler) closeRoom(roomID pgtype.UUID, reason string) {
	h.hub.CloseRoom(roomID.String(), reason)
}

func (h *RoomHandler) closeUser(roomID pgtype.UUID, userID string, reason string) {
	h.hub.CloseUser(roomID.String(), userID, reason)
}

func (h *RoomHandler) handleAutomaticLeave(result roomDTO.RoomMutationResult[roomDTO.LeaveRoomResult]) {
	h.broadcastRoomEvents(result.Events)
	if result.Value.DeletedRoomID != nil {
		h.presence.ForgetRoom(*result.Value.DeletedRoomID)
		h.closeRoom(*result.Value.DeletedRoomID, "room deleted after disconnect")
	}
}
