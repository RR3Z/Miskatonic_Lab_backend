package room

import (
	"net/http"

	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room/helpers"
	roomDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/coder/websocket"
)

func (h *RoomHandler) serveRoomWS(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	if err := h.service.TouchRoomActivity(r.Context(), roomDTO.TouchRoomActivityInput{
		RoomID: roomID,
		UserID: userID,
	}); err != nil {
		http.Error(w, "room websocket unavailable", http.StatusForbidden)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(h.hub, h.dispatcher, roomID, userID, conn)
	h.hub.Register <- client

	go client.WriteLoop(r.Context())
	client.ReadLoop(r.Context())
}
