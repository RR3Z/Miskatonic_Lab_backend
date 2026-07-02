package ws_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgtype"
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

type roomCommand struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type fakeRoomEventService struct {
	inputs chan roomModel.CreateChatMessageInput
}

func newFakeRoomEventService() *fakeRoomEventService {
	return &fakeRoomEventService{
		inputs: make(chan roomModel.CreateChatMessageInput, 8),
	}
}

func (f *fakeRoomEventService) CreateChatMessage(_ context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error) {
	f.inputs <- input

	payload, err := json.Marshal(roomEvents.ChatMessagePayload{Text: input.Text})
	if err != nil {
		return roomModel.RoomEventModel{}, err
	}

	return roomModel.RoomEventModel{
		RoomID:  input.RoomID,
		ActorID: input.ActorID,
		Type:    string(roomEvents.EventChatMessage),
		Payload: payload,
	}, nil
}

func (f *fakeRoomEventService) waitForCreateChatInput(t *testing.T, ctx context.Context) roomModel.CreateChatMessageInput {
	t.Helper()

	select {
	case input := <-f.inputs:
		return input
	case <-ctx.Done():
		t.Fatal("timed out waiting for chat message persistence")
		return roomModel.CreateChatMessageInput{}
	}
}

func startRoomWSServer(t *testing.T, hub *ws.RoomHub, service ws.RoomEventService, roomID pgtype.UUID) (<-chan string, string) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	registered := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		client := ws.NewClient(hub, service, roomID, userID, conn)
		hub.Register <- client
		registered <- userID

		go client.WriteLoop(r.Context())
		client.ReadLoop(r.Context())
	}))
	t.Cleanup(server.Close)

	return registered, "ws" + strings.TrimPrefix(server.URL, "http")
}

func dialRoomClient(t *testing.T, ctx context.Context, serverURL string, userID string) *websocket.Conn {
	t.Helper()

	conn, _, err := websocket.Dial(ctx, serverURL+"/ws?user="+userID, nil)
	require.NoError(t, err)
	return conn
}

func waitForRegisteredUsers(t *testing.T, ctx context.Context, registered <-chan string, expected ...string) {
	t.Helper()

	want := make(map[string]struct{}, len(expected))
	for _, userID := range expected {
		want[userID] = struct{}{}
	}

	var mu sync.Mutex
	seen := make(map[string]struct{}, len(expected))
	for len(seen) < len(want) {
		select {
		case userID := <-registered:
			mu.Lock()
			if _, ok := want[userID]; ok {
				seen[userID] = struct{}{}
			}
			mu.Unlock()
		case <-ctx.Done():
			t.Fatalf("timed out waiting for registered websocket clients: got %v, want %v", seen, want)
		}
	}
}

func testWSUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
