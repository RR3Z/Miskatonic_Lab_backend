package ws_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	wsCommands "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/commands"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/coder/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func startRoomWSServer(t *testing.T, hub *ws.RoomHub, service wsCommands.RoomEventService, roomID pgtype.UUID) (<-chan string, string) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)
	dispatcher := wsCommands.NewCommandDispatcher(service)

	registered := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		client := ws.NewClient(hub, dispatcher, roomID, userID, conn)
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
