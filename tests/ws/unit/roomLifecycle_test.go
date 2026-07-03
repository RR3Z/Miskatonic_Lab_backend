package ws_test

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	"github.com/coder/websocket"
)

func TestRoomHubCloseRoomClosesActiveWebSocketConnection(t *testing.T) {
	roomID := testWSUUID("11111111-1111-1111-1111-111111111111")
	service := newFakeRoomEventService()
	hub := ws.NewRoomHub()
	registered, serverURL := startRoomWSServer(t, hub, service, roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn := dialRoomClient(t, ctx, serverURL, "user_1")
	defer conn.Close(websocket.StatusNormalClosure, "test done")
	waitForRegisteredUsers(t, ctx, registered, "user_1")

	hub.CloseRoom(roomID.String(), "room deleted")

	requireWebSocketClosed(t, ctx, conn)
}

func TestRoomHubCloseRoomOnlyClosesTargetRoomAndIsRepeatSafe(t *testing.T) {
	roomAID := "11111111-1111-1111-1111-111111111111"
	roomBID := "22222222-2222-2222-2222-222222222222"
	hub := ws.NewRoomHub()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	roomAClient, roomAEvents := ws.NewTestClientWithUser(hub, roomAID, "user_a")
	roomBClient, roomBEvents := ws.NewTestClientWithUser(hub, roomBID, "user_b")
	hub.Register <- roomAClient
	hub.Register <- roomBClient

	hub.CloseRoom(roomAID, "room deleted")
	hub.CloseRoom(roomAID, "room deleted again")

	requireEventChannelClosed(t, roomAEvents)
	requireEventChannelOpen(t, roomBEvents)
}
