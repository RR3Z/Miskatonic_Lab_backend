package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/stretchr/testify/require"
)

func TestRoomRealtimeWorkabilityScenario(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	gm := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	room, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:  gm.ID,
		Password: "keeper-password",
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:      room.ID,
		UserID:      firstPlayer.ID,
		InviteToken: room.InviteToken,
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:   room.ID,
		UserID:   secondPlayer.ID,
		Password: "keeper-password",
	})
	require.NoError(t, err)

	firstCharacter := createRoomTestCharacter(t, subject, firstPlayer.ID)
	secondCharacter := createRoomTestCharacter(t, subject, secondPlayer.ID)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      firstPlayer.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      secondPlayer.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)

	gmCharacters, err := service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: gm.ID,
	})
	require.NoError(t, err)
	requireSelectedCharacterUsers(t, gmCharacters, firstPlayer.ID, secondPlayer.ID)

	firstPlayerCharacters, err := service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: firstPlayer.ID,
	})
	require.NoError(t, err)
	require.Len(t, firstPlayerCharacters, 1)
	require.Equal(t, firstPlayer.ID, firstPlayerCharacters[0].UserID)
	require.Equal(t, firstCharacter.ID, firstPlayerCharacters[0].Character.ID)

	_, err = service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: gm.ID,
		Text:    "the seance begins",
	})
	require.NoError(t, err)
	time.Sleep(5 * time.Millisecond)

	_, err = service.CreateDiceRollRoomEvent(context.Background(), model.CreateDiceRollRoomEventInput{
		RoomID:      room.ID,
		ActorID:     firstPlayer.ID,
		RollID:      "workability-roll-1",
		CharacterID: firstCharacter.ID.String(),
		Expression:  "1d20",
		Result:      13,
		Details:     []byte(`[{"type":"dice","sides":20,"rolls":[13]}]`),
	})
	require.NoError(t, err)
	time.Sleep(5 * time.Millisecond)

	sourceEvent := "character.health.upsert_succeeded"
	_, err = service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		CharacterID: firstCharacter.ID,
		ActorID:     firstPlayer.ID,
		Change: model.CharacterChangedRoomEventChange{
			Resource:    "health",
			Action:      "upsert",
			SourceEvent: &sourceEvent,
		},
	})
	require.NoError(t, err)
	time.Sleep(5 * time.Millisecond)

	_, err = service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		CharacterID: secondCharacter.ID,
		ActorID:     secondPlayer.ID,
		Change: model.CharacterChangedRoomEventChange{
			Resource:    "skill",
			Action:      "update",
			SourceEvent: &sourceEvent,
		},
	})
	require.NoError(t, err)

	gmEvents, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: gm.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, gmEvents, 4)
	require.Equal(t, []string{
		string(roomEvents.EventChatMessage),
		string(roomEvents.EventDiceRoll),
		string(roomEvents.EventCharacterChanged),
		string(roomEvents.EventCharacterChanged),
	}, roomEventTypes(gmEvents))
	require.ElementsMatch(t, []string{firstCharacter.ID.String(), secondCharacter.ID.String()}, characterChangedCharacterIDs(t, gmEvents))

	firstPlayerEvents, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: firstPlayer.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, firstPlayerEvents, 3)
	require.Equal(t, []string{
		string(roomEvents.EventChatMessage),
		string(roomEvents.EventDiceRoll),
		string(roomEvents.EventCharacterChanged),
	}, roomEventTypes(firstPlayerEvents))
	require.Equal(t, []string{firstCharacter.ID.String()}, characterChangedCharacterIDs(t, firstPlayerEvents))
	requireDiceRollPayload(t, firstPlayerEvents[1], "workability-roll-1", firstCharacter.ID.String(), "1d20", int32(13))

	secondPlayerEvents, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: secondPlayer.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, secondPlayerEvents, 3)
	require.Equal(t, []string{secondCharacter.ID.String()}, characterChangedCharacterIDs(t, secondPlayerEvents))
}

func roomEventTypes(events []model.RoomEventModel) []string {
	types := make([]string, 0, len(events))
	for _, event := range events {
		types = append(types, event.Type)
	}
	return types
}

func requireDiceRollPayload(
	t *testing.T,
	event model.RoomEventModel,
	rollID string,
	characterID string,
	expression string,
	result int32,
) {
	t.Helper()

	require.Equal(t, string(roomEvents.EventDiceRoll), event.Type)

	var payload roomEvents.DiceRollPayload
	require.NoError(t, json.Unmarshal(event.Payload, &payload))
	require.Equal(t, rollID, payload.RollID)
	require.Equal(t, characterID, payload.CharacterID)
	require.Equal(t, expression, payload.Expression)
	require.Equal(t, result, payload.Result)
}
