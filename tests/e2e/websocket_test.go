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

func TestE2ERoomWebSocketCursorRecoversEventsAfterReconnect(t *testing.T) {
	owner := newE2ESubject(t)
	player := newSecondE2ESubject(t)
	password := "e2e-recovery-" + e2eHash(owner.userID+player.userID)
	room := owner.createRoom(t, password)
	t.Cleanup(func() {
		owner.deleteRoom(t, room.ID)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ownerConn := dialE2ERoomSocket(t, ctx, owner, room.ID)
	player.joinRoom(t, room.ID, password)
	playerConn := dialE2ERoomSocket(t, ctx, player, room.ID)
	t.Cleanup(func() {
		ownerConn.Close(websocket.StatusNormalClosure, "test done")
		playerConn.Close(websocket.StatusNormalClosure, "test done")
	})

	firstText := "before reconnect " + e2eHash(time.Now().UTC().String())
	require.NoError(t, wsjson.Write(ctx, playerConn, e2eRoomCommand{
		Type:    "chat.message",
		Payload: e2eChatPayload{Text: firstText},
	}))
	firstEvent := readE2ERoomSocketEvent(t, ctx, ownerConn, "chat.message")
	require.Greater(t, firstEvent.Sequence, int64(0))
	requireChatPayloadText(t, firstEvent.Payload, firstText)

	require.NoError(t, ownerConn.Close(websocket.StatusNormalClosure, "simulate disconnected tab"))
	secondText := "while disconnected " + e2eHash(time.Now().UTC().Add(time.Second).String())
	require.NoError(t, wsjson.Write(ctx, playerConn, e2eRoomCommand{
		Type:    "chat.message",
		Payload: e2eChatPayload{Text: secondText},
	}))

	ownerReconnect := dialE2ERoomSocket(t, ctx, owner, room.ID)
	defer ownerReconnect.Close(websocket.StatusNormalClosure, "test done")
	recoveredEvents := owner.listRoomEventsAfter(t, room.ID, firstEvent.Sequence)
	require.NotEmpty(t, recoveredEvents)

	foundRecoveredChat := false
	for _, event := range recoveredEvents {
		require.Greater(t, event.Sequence, firstEvent.Sequence)
		if event.Type != "chat.message" {
			continue
		}
		var chat e2eChatPayload
		require.NoError(t, json.Unmarshal(event.Payload, &chat))
		if chat.Text == secondText {
			foundRecoveredChat = true
		}
	}
	require.True(t, foundRecoveredChat, "cursor response must contain the message sent during disconnect")
}

func dialE2ERoomSocket(t *testing.T, ctx context.Context, subject *e2eSubject, roomID string) *websocket.Conn {
	t.Helper()

	conn, _, err := websocket.Dial(ctx, subject.wsURL(t, "/api/rooms/"+url.PathEscape(roomID)+"/ws"), &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": []string{subject.authorization(t)}},
	})
	require.NoError(t, err)
	return conn
}

func readE2ERoomSocketEvent(t *testing.T, ctx context.Context, conn *websocket.Conn, eventType string) e2eRoomSocketEvent {
	t.Helper()

	for {
		var event e2eRoomSocketEvent
		require.NoError(t, wsjson.Read(ctx, conn, &event))
		if event.Type == eventType {
			return event
		}
	}
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
