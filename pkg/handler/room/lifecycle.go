package room

import "github.com/jackc/pgx/v5/pgtype"

func (h *RoomHandler) CloseDeletedRooms(roomIDs []pgtype.UUID, reason string) {
	for _, roomID := range roomIDs {
		h.closeRoom(roomID, reason)
	}
}

func (h *RoomHandler) closeRoom(roomID pgtype.UUID, reason string) {
	h.hub.CloseRoom(roomID.String(), reason)
}
