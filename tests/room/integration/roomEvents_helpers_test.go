package tests

import (
	"context"
	"encoding/json"
	"testing"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func requireRoomCharacterChangedEvent(
	t *testing.T,
	subject *roomIntegrationSubject,
	roomID pgtype.UUID,
	userID string,
	characterID string,
	resource string,
	action string,
	resourceID *string,
	sourceEvent *string,
) {
	t.Helper()

	events, err := subject.queries.ListRoomEvents(context.Background(), db.ListRoomEventsParams{
		RoomID:     roomID,
		UserID:     userID,
		LimitCount: 10,
	})
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, string(roomEvents.EventCharacterChanged), events[0].EventType)

	var payload roomEvents.CharacterChangedPayload
	require.NoError(t, json.Unmarshal(events[0].Payload, &payload))
	require.Equal(t, characterID, payload.CharacterID)
	require.Equal(t, resource, payload.Resource)
	require.Equal(t, action, payload.Action)
	require.Equal(t, resourceID, payload.ResourceID)
	require.Equal(t, sourceEvent, payload.SourceEvent)
}

func requireRoomEventTypes(t *testing.T, events []model.RoomEventModel, expectedTypes ...string) {
	t.Helper()

	actualTypes := make([]string, 0, len(events))
	for _, event := range events {
		actualTypes = append(actualTypes, event.Type)
	}

	for _, expectedType := range expectedTypes {
		require.Contains(t, actualTypes, expectedType)
	}
}

func characterChangedCharacterIDs(t *testing.T, events []model.RoomEventModel) []string {
	t.Helper()

	characterIDs := make([]string, 0)
	for _, event := range events {
		if event.Type != string(roomEvents.EventCharacterChanged) {
			continue
		}

		var payload roomEvents.CharacterChangedPayload
		require.NoError(t, json.Unmarshal(event.Payload, &payload))
		characterIDs = append(characterIDs, payload.CharacterID)
	}

	return characterIDs
}
