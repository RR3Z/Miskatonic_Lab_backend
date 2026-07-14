package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/require"
)

func TestE2ERoomWebSocketChatPersistsAndBroadcasts(t *testing.T) {
	subject := newE2ESubject(t)

	room := subject.createRoom(t, "e2e-ws-"+e2eHash(subject.userID))
	require.Equal(t, subject.userID, room.OwnerID)
	t.Cleanup(func() {
		subject.deleteRoom(t, room.ID)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, subject.wsURL(t, "/api/rooms/"+url.PathEscape(room.ID)+"/ws"), &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": []string{subject.authorization(t)}},
	})
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	text := "hello from live ws " + e2eHash(time.Now().UTC().String())
	err = wsjson.Write(ctx, conn, e2eRoomCommand{
		Type:    "chat.message",
		Payload: e2eChatPayload{Text: text},
	})
	require.NoError(t, err)

	var socketEvent e2eRoomSocketEvent
	require.NoError(t, wsjson.Read(ctx, conn, &socketEvent))
	require.Equal(t, "chat.message", socketEvent.Type)
	require.Equal(t, room.ID, socketEvent.RoomID)
	require.Equal(t, subject.userID, socketEvent.ActorID)
	requireChatPayloadText(t, socketEvent.Payload, text)

	events := subject.waitForRoomEvents(t, room.ID, "chat.message")
	require.NotEmpty(t, events)
	requireRoomHistoryChatText(t, events, text)
}

func requireChatPayloadText(t *testing.T, payload json.RawMessage, expected string) {
	t.Helper()

	var chat e2eChatPayload
	require.NoError(t, json.Unmarshal(payload, &chat))
	require.Equal(t, expected, chat.Text)
}

func requireRoomHistoryChatText(t *testing.T, events []e2eRoomEventResponse, expected string) {
	t.Helper()

	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Type != "chat.message" {
			continue
		}
		requireChatPayloadText(t, events[i].Payload, expected)
		return
	}
	t.Fatalf("room history does not contain chat.message with text %q", expected)
}
