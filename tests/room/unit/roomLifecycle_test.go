package tests

import (
	"net/http"
	"testing"

	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestDeleteRoomClosesActiveRoomClient(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	hub := ws.NewRoomHub()
	events := registerRoomUnitTestClient(t, hub, roomID)
	roomService := &fakeRoomHandlerService{}
	router := newRoomHandlerTestRouterWithHub(roomService, hub)

	recorder := performRoomRequest(router, http.MethodDelete, "/api/rooms/11111111-1111-1111-1111-111111111111/", "")

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, roomService.deleteCalls)
	requireRoomUnitClientClosed(t, events)
}

func TestLastMemberLeaveClosesActiveRoomClient(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	hub := ws.NewRoomHub()
	events := registerRoomUnitTestClient(t, hub, roomID)
	deletedRoomID := roomID
	roomService := &fakeRoomHandlerService{
		leaveResult: roomModels.LeaveRoomResult{DeletedRoomID: &deletedRoomID},
	}
	router := newRoomHandlerTestRouterWithHub(roomService, hub)

	recorder := performRoomRequest(router, http.MethodDelete, "/api/rooms/11111111-1111-1111-1111-111111111111/leave", "")

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, roomService.leaveCalls)
	requireRoomUnitClientClosed(t, events)
}

func TestCleanupDeletedRoomsCloseActiveRoomClient(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	hub := ws.NewRoomHub()
	events := registerRoomUnitTestClient(t, hub, roomID)
	handler := roomHandler.NewWithHub(&fakeRoomHandlerService{}, hub)

	handler.CloseDeletedRooms([]pgtype.UUID{roomID}, "room deleted by cleanup")

	requireRoomUnitClientClosed(t, events)
}
