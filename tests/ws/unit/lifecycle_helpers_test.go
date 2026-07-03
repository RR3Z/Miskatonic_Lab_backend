package ws_test

import (
	"context"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/require"
)

func requireWebSocketClosed(t *testing.T, ctx context.Context, conn *websocket.Conn) {
	t.Helper()

	var event any
	err := wsjson.Read(ctx, conn, &event)
	require.Error(t, err)
}

func requireNoWebSocketEvent(t *testing.T, conn *websocket.Conn) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var event any
	err := wsjson.Read(ctx, conn, &event)
	require.Error(t, err)
}

func requireEventChannelClosed(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	require.Eventually(t, func() bool {
		select {
		case _, ok := <-events:
			return !ok
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)
}

func requireEventChannelOpen(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	select {
	case _, ok := <-events:
		require.True(t, ok)
	default:
	}
}

func requireNoRoomHubEvent(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	select {
	case event := <-events:
		t.Fatalf("unexpected room hub event: %#v", event)
	case <-time.After(100 * time.Millisecond):
	}
}
