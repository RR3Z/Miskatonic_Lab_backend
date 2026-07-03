package ws_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/require"
)

func TestRoomWebsocketChatMessageIsSavedBeforeBroadcast(t *testing.T) {
	roomID := testWSUUID("11111111-1111-1111-1111-111111111111")
	service := newFakeRoomEventService()
	hub := ws.NewRoomHub()
	_, serverURL := startRoomWSServer(t, hub, service, roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn := dialRoomClient(t, ctx, serverURL, "user_1")
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	err := wsjson.Write(ctx, conn, roomCommand{
		Type:    string(roomEvents.EventChatMessage),
		Payload: json.RawMessage(`{"text":"hello","actor_id":"spoofed"}`),
	})
	require.NoError(t, err)

	input := service.waitForCreateChatInput(t, ctx)
	require.Equal(t, roomID, input.RoomID)
	require.Equal(t, "user_1", input.ActorID)
	require.Equal(t, "hello", input.Text)

	var event roomEvents.Event
	require.NoError(t, wsjson.Read(ctx, conn, &event))
	require.Equal(t, string(roomEvents.EventChatMessage), event.Type)
	require.Equal(t, roomID.String(), event.RoomID)
	require.Equal(t, "user_1", event.ActorID)
	require.Equal(t, map[string]any{"text": "hello"}, event.Payload)
}

func TestRoomHubBroadcastReachesFastClientWhenSlowClientIsConnected(t *testing.T) {
	roomID := testWSUUID("11111111-1111-1111-1111-111111111111")
	service := newFakeRoomEventService()
	hub := ws.NewRoomHub()
	registered, serverURL := startRoomWSServer(t, hub, service, roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	slowConn := dialRoomClient(t, ctx, serverURL, "slow")
	defer slowConn.Close(websocket.StatusNormalClosure, "test done")
	fastConn := dialRoomClient(t, ctx, serverURL, "fast")
	defer fastConn.Close(websocket.StatusNormalClosure, "test done")

	waitForRegisteredUsers(t, ctx, registered, "slow", "fast")

	hub.Broadcast(roomEvents.Event{
		Type:    string(roomEvents.EventChatMessage),
		RoomID:  roomID.String(),
		ActorID: "system",
		Payload: roomEvents.ChatMessagePayload{Text: "hello"},
	})

	var event roomEvents.Event
	require.NoError(t, wsjson.Read(ctx, fastConn, &event))
	require.Equal(t, string(roomEvents.EventChatMessage), event.Type)
	require.Equal(t, "system", event.ActorID)
	require.Equal(t, map[string]any{"text": "hello"}, event.Payload)
}

func TestRoomWebsocketReturnsCommandErrorForUnknownCommandType(t *testing.T) {
	roomID := testWSUUID("11111111-1111-1111-1111-111111111111")
	service := newFakeRoomEventService()
	hub := ws.NewRoomHub()
	_, serverURL := startRoomWSServer(t, hub, service, roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn := dialRoomClient(t, ctx, serverURL, "user_1")
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	err := wsjson.Write(ctx, conn, roomCommand{
		Type:    "dice.roll",
		Payload: json.RawMessage(`{"expression":"1d20"}`),
	})
	require.NoError(t, err)

	var event roomEvents.Event
	require.NoError(t, wsjson.Read(ctx, conn, &event))
	require.Equal(t, string(roomEvents.EventCommandError), event.Type)
	require.Equal(t, roomID.String(), event.RoomID)
	require.Equal(t, "user_1", event.ActorID)
	require.Equal(t, map[string]any{
		"code":    "common.invalid_request",
		"message": "unsupported room command type",
		"details": []any{
			map[string]any{
				"type":   "validation",
				"target": "command.type",
				"reason": "unsupported",
			},
		},
	}, event.Payload)
	require.Len(t, service.inputs, 0)

	err = wsjson.Write(ctx, conn, roomCommand{
		Type:    string(roomEvents.EventChatMessage),
		Payload: json.RawMessage(`{"text":"still connected"}`),
	})
	require.NoError(t, err)

	input := service.waitForCreateChatInput(t, ctx)
	require.Equal(t, "still connected", input.Text)

	require.NoError(t, wsjson.Read(ctx, conn, &event))
	require.Equal(t, string(roomEvents.EventChatMessage), event.Type)
	require.Equal(t, map[string]any{"text": "still connected"}, event.Payload)
}

func TestRoomWebsocketReturnsCommandErrorForMalformedChatPayload(t *testing.T) {
	roomID := testWSUUID("11111111-1111-1111-1111-111111111111")
	service := newFakeRoomEventService()
	hub := ws.NewRoomHub()
	_, serverURL := startRoomWSServer(t, hub, service, roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn := dialRoomClient(t, ctx, serverURL, "user_1")
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	err := wsjson.Write(ctx, conn, roomCommand{
		Type:    string(roomEvents.EventChatMessage),
		Payload: json.RawMessage(`"not an object"`),
	})
	require.NoError(t, err)

	var event roomEvents.Event
	require.NoError(t, wsjson.Read(ctx, conn, &event))
	require.Equal(t, string(roomEvents.EventCommandError), event.Type)
	require.Equal(t, roomID.String(), event.RoomID)
	require.Equal(t, "user_1", event.ActorID)
	require.Equal(t, map[string]any{
		"code":    "common.invalid_request",
		"message": "invalid room command payload",
		"details": []any{
			map[string]any{
				"type":   "parse",
				"target": "command.payload",
				"reason": "invalid_format",
			},
		},
	}, event.Payload)
	require.Len(t, service.inputs, 0)

	err = wsjson.Write(ctx, conn, roomCommand{
		Type:    string(roomEvents.EventChatMessage),
		Payload: json.RawMessage(`{"text":"valid after error"}`),
	})
	require.NoError(t, err)

	input := service.waitForCreateChatInput(t, ctx)
	require.Equal(t, "valid after error", input.Text)
}

type roomCommand struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
