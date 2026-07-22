package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func TestRoomJoinWithoutWebSocketAutoLeavesAfterGrace(t *testing.T) {
	service := &fakeRoomHandlerService{
		member: roomModels.RoomMemberModel{Role: "player"},
	}
	hub := ws.NewRoomHub()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)
	router := newRoomHandlerTestRouterWithHubAndPresence(service, hub, roomHandler.PresenceConfig{
		DisconnectGrace: 20 * time.Millisecond,
	})

	recorder := performRoomRequest(
		router,
		http.MethodPost,
		"/api/rooms/11111111-1111-1111-1111-111111111111/join",
		`{"invite_token":"token"}`,
	)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Eventually(t, func() bool {
		return service.leaveCalls == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, testRoomUnitUUID("11111111-1111-1111-1111-111111111111"), service.leaveInput.RoomID)
	require.Equal(t, "user_1", service.leaveInput.UserID)
}

func TestRoomPresenceWaitsForUsersLastWebSocket(t *testing.T) {
	const grace = 50 * time.Millisecond
	service := &fakeRoomHandlerService{
		member: roomModels.RoomMemberModel{Role: "player"},
	}
	hub := ws.NewRoomHub()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)
	router := newRoomHandlerTestRouterWithHubAndPresence(service, hub, roomHandler.PresenceConfig{
		DisconnectGrace: grace,
	})
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	roomID := "11111111-1111-1111-1111-111111111111"
	joinRequest, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/rooms/"+roomID+"/join",
		bytes.NewBufferString(`{"invite_token":"token"}`),
	)
	require.NoError(t, err)
	joinRequest.Header.Set("Content-Type", "application/json")
	joinResponse, err := server.Client().Do(joinRequest)
	require.NoError(t, err)
	joinResponse.Body.Close()
	require.Equal(t, http.StatusOK, joinResponse.StatusCode)

	websocketURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/rooms/" + roomID + "/ws"
	firstConnection, _, err := websocket.Dial(ctx, websocketURL, nil)
	require.NoError(t, err)
	t.Cleanup(func() { firstConnection.Close(websocket.StatusNormalClosure, "test done") })
	secondConnection, _, err := websocket.Dial(ctx, websocketURL, nil)
	require.NoError(t, err)
	t.Cleanup(func() { secondConnection.Close(websocket.StatusNormalClosure, "test done") })

	time.Sleep(2 * grace)
	require.Zero(t, service.leaveCalls)

	require.NoError(t, firstConnection.Close(websocket.StatusNormalClosure, "first tab closed"))
	time.Sleep(2 * grace)
	require.Zero(t, service.leaveCalls)

	require.NoError(t, secondConnection.Close(websocket.StatusNormalClosure, "last tab closed"))
	require.Eventually(t, func() bool {
		return service.leaveCalls == 1
	}, time.Second, 10*time.Millisecond)
}
