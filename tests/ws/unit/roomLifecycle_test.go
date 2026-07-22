package ws_test

import (
	"context"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
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

func TestRoomHubCloseUserClosesOnlyNamedUsersConnections(t *testing.T) {
	roomID := "11111111-1111-1111-1111-111111111111"
	hub := ws.NewRoomHub()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	firstTargetClient, firstTargetEvents := ws.NewTestClientWithUser(hub, roomID, "target")
	secondTargetClient, secondTargetEvents := ws.NewTestClientWithUser(hub, roomID, "target")
	otherClient, otherEvents := ws.NewTestClientWithUser(hub, roomID, "other")
	hub.Register <- firstTargetClient
	hub.Register <- secondTargetClient
	hub.Register <- otherClient

	hub.CloseUser(roomID, "target", "removed")

	requireEventChannelClosed(t, firstTargetEvents)
	requireEventChannelClosed(t, secondTargetEvents)
	requireEventChannelOpen(t, otherEvents)
}

func TestRoomHubSendToUsersTargetsOnlyNamedUsers(t *testing.T) {
	roomID := "11111111-1111-1111-1111-111111111111"
	hub := ws.NewRoomHub()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	targetClient, targetEvents := ws.NewTestClientWithUser(hub, roomID, "target")
	otherClient, otherEvents := ws.NewTestClientWithUser(hub, roomID, "other")
	hub.Register <- targetClient
	hub.Register <- otherClient

	hub.SendToUsers(roomID, []string{"target"}, roomEvents.Event{
		Type:    string(roomEvents.EventCharacterChanged),
		RoomID:  roomID,
		ActorID: "system",
		Payload: roomEvents.CharacterChangedPayload{CharacterID: "character_1", Resource: "health", Action: "upsert"},
	})

	select {
	case event := <-targetEvents:
		require.Equal(t, string(roomEvents.EventCharacterChanged), event.Type)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for targeted event")
	}
	requireNoRoomHubEvent(t, otherEvents)
}

func TestRoomHubSendToUsersWithEmptyTargetsNoOps(t *testing.T) {
	roomID := "11111111-1111-1111-1111-111111111111"
	hub := ws.NewRoomHub()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	client, events := ws.NewTestClientWithUser(hub, roomID, "user_1")
	hub.Register <- client

	hub.SendToUsers(roomID, nil, roomEvents.Event{
		Type:    string(roomEvents.EventCharacterChanged),
		RoomID:  roomID,
		ActorID: "system",
		Payload: roomEvents.CharacterChangedPayload{CharacterID: "character_1", Resource: "health", Action: "upsert"},
	})

	requireNoRoomHubEvent(t, events)
}
