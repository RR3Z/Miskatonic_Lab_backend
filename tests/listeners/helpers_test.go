package tests

import (
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/stretchr/testify/require"
)

func requireCharacterChangedRealtimeEvent(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	select {
	case event := <-events:
		require.Equal(t, string(roomEvents.EventCharacterChanged), event.Type)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for character.changed realtime event")
	}
}

func requireNoRealtimeEvent(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	select {
	case event := <-events:
		t.Fatalf("unexpected realtime event: %#v", event)
	case <-time.After(100 * time.Millisecond):
	}
}
