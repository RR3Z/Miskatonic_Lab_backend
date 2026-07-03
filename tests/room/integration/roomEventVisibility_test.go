package tests

import (
	"context"
	"testing"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/stretchr/testify/require"
)

func TestRoomServiceListRoomEventsFiltersCharacterChangedByRole(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	gm := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, gm.ID)
	addRoomTestMember(t, subject, room.ID, gm.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, firstPlayer.ID, roomService.ROLE_PLAYER)
	addRoomTestMember(t, subject, room.ID, secondPlayer.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	firstCharacter := createRoomTestCharacter(t, subject, firstPlayer.ID)
	secondCharacter := createRoomTestCharacter(t, subject, secondPlayer.ID)
	_, err := service.SelectCharacter(context.Background(), model.SelectCharacterInput{
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

	_, err = service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: gm.ID,
		Text:    "everyone can read this",
	})
	require.NoError(t, err)

	sourceEvent := "character.health.upsert_succeeded"
	firstChangedEvents, err := service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		CharacterID: firstCharacter.ID,
		ActorID:     firstPlayer.ID,
		Change: model.CharacterChangedRoomEventChange{
			Resource:    "health",
			Action:      "upsert",
			SourceEvent: &sourceEvent,
		},
	})
	require.NoError(t, err)
	require.Len(t, firstChangedEvents, 1)
	require.ElementsMatch(t, []string{firstPlayer.ID, gm.ID}, firstChangedEvents[0].TargetUserIDs)

	_, err = service.CreateDiceRollRoomEvent(context.Background(), model.CreateDiceRollRoomEventInput{
		RoomID:      room.ID,
		ActorID:     secondPlayer.ID,
		RollID:      "roll_1",
		CharacterID: secondCharacter.ID.String(),
		Expression:  "1d20",
		Result:      12,
		Details:     []byte(`[]`),
	})
	require.NoError(t, err)

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
	require.ElementsMatch(t, []string{firstCharacter.ID.String(), secondCharacter.ID.String()}, characterChangedCharacterIDs(t, gmEvents))

	firstPlayerEvents, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: firstPlayer.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, firstPlayerEvents, 3)
	requireRoomEventTypes(t, firstPlayerEvents, string(roomEvents.EventChatMessage), string(roomEvents.EventDiceRoll))
	require.Equal(t, []string{firstCharacter.ID.String()}, characterChangedCharacterIDs(t, firstPlayerEvents))

	secondPlayerEvents, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: secondPlayer.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, secondPlayerEvents, 3)
	requireRoomEventTypes(t, secondPlayerEvents, string(roomEvents.EventChatMessage), string(roomEvents.EventDiceRoll))
	require.Equal(t, []string{secondCharacter.ID.String()}, characterChangedCharacterIDs(t, secondPlayerEvents))
}
