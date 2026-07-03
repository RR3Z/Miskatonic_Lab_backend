package tests

import (
	"context"
	"encoding/json"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

func TestRoomServiceCreateRoomCreatesOwnerMemberAndInviteToken(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	maxPlayers := int32(3)

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:    owner.ID,
		MaxPlayers: &maxPlayers,
		Password:   "keeper-password",
	})
	require.NoError(t, err)
	require.True(t, roomModel.ID.Valid)
	require.Equal(t, owner.ID, roomModel.OwnerID)
	require.Equal(t, int32(3), roomModel.MaxPlayers)
	require.NotEmpty(t, roomModel.InviteToken)
	require.Len(t, roomModel.Members, 1)
	require.Equal(t, owner.ID, roomModel.Members[0].UserID)
	require.Equal(t, "gm", roomModel.Members[0].Role)

	member, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: roomModel.ID, UserID: owner.ID})
	require.NoError(t, err)
	require.Equal(t, "gm", member.Role)

	meta, err := subject.queries.GetRoomJoinMetaData(context.Background(), roomModel.ID)
	require.NoError(t, err)
	require.NotEqual(t, "keeper-password", meta.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(meta.PasswordHash), []byte("keeper-password")))
}

func TestRoomServiceCreateRoomDefaultsAndValidatesMaxPlayers(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{OwnerID: owner.ID, Password: "keeper-password"})
	require.NoError(t, err)
	require.Equal(t, roomService.DEFAULT_MAX_PLAYERS, roomModel.MaxPlayers)

	invalidMaxPlayers := int32(0)
	_, err = service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:    owner.ID,
		MaxPlayers: &invalidMaxPlayers,
		Password:   "keeper-password",
	})
	require.ErrorIs(t, err, roomService.ErrInvalidInput)

	_, err = service.CreateRoom(context.Background(), model.CreateRoomInput{OwnerID: owner.ID})
	require.ErrorIs(t, err, roomService.ErrInvalidPassword)

	_, err = service.CreateRoom(context.Background(), model.CreateRoomInput{OwnerID: owner.ID, Password: "   "})
	require.ErrorIs(t, err, roomService.ErrInvalidPassword)
}

func TestRoomServiceJoinRoomRequiresInviteOrPasswordAndCapacity(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	thirdPlayer := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	maxPlayers := int32(2)
	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:    owner.ID,
		MaxPlayers: &maxPlayers,
		Password:   "keeper-password",
	})
	require.NoError(t, err)
	room := db.Room{ID: roomModel.ID, InviteToken: roomModel.InviteToken}

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, InviteToken: "wrong_token", UserID: firstPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomNotFound)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, Password: "wrong-password", UserID: firstPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomNotFound)

	memberModel, err := service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, InviteToken: room.InviteToken, UserID: firstPlayer.ID},
	)
	require.NoError(t, err)
	require.Equal(t, firstPlayer.ID, memberModel.UserID)
	require.Equal(t, "player", memberModel.Role)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, InviteToken: room.InviteToken, UserID: firstPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrAlreadyMember)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, InviteToken: room.InviteToken, UserID: secondPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomFull)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, Password: "keeper-password", UserID: thirdPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomFull)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: room.ID, UserID: thirdPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrInvalidInput)
}

func TestRoomServiceJoinRoomWithPassword(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	player := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:  owner.ID,
		Password: "keeper-password",
	})
	require.NoError(t, err)

	memberModel, err := service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: roomModel.ID, Password: "keeper-password", UserID: player.ID},
	)
	require.NoError(t, err)
	require.Equal(t, player.ID, memberModel.UserID)
	require.Equal(t, roomService.ROLE_PLAYER, memberModel.Role)
}

func TestRoomServiceJoinRoomAcceptsEitherValidInviteOrPassword(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:  owner.ID,
		Password: "keeper-password",
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:      roomModel.ID,
		UserID:      firstPlayer.ID,
		InviteToken: "wrong-invite",
		Password:    "keeper-password",
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:      roomModel.ID,
		UserID:      secondPlayer.ID,
		InviteToken: roomModel.InviteToken,
		Password:    "wrong-password",
	})
	require.NoError(t, err)
}

func TestRoomServiceTransferOwnershipDoesNotChangeRole(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	transferred, err := service.TransferOwnership(context.Background(), model.TransferOwnershipInput{
		RoomID:     room.ID,
		OwnerID:    owner.ID,
		NewOwnerID: memberUser.ID,
	})
	require.NoError(t, err)
	require.Equal(t, memberUser.ID, transferred.OwnerID)

	member, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, "player", member.Role)

	_, err = service.TransferOwnership(context.Background(), model.TransferOwnershipInput{
		RoomID:     room.ID,
		OwnerID:    memberUser.ID,
		NewOwnerID: outsider.ID,
	})
	require.ErrorIs(t, err, roomService.ErrNotOwner)

	_, err = service.TransferOwnership(context.Background(), model.TransferOwnershipInput{
		RoomID:  room.ID,
		OwnerID: owner.ID,
	})
	require.ErrorIs(t, err, roomService.ErrInvalidInput)
}

func TestRoomServiceNonOwnerUpdateAndDeletePreserveRoom(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:    owner.ID,
		MaxPlayers: nil,
		Password:   "keeper-password",
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:      roomModel.ID,
		UserID:      memberUser.ID,
		InviteToken: roomModel.InviteToken,
	})
	require.NoError(t, err)

	_, err = service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:     roomModel.ID,
		OwnerID:    memberUser.ID,
		MaxPlayers: 1,
	})
	require.ErrorIs(t, err, roomService.ErrNotOwner)

	err = service.DeleteRoom(context.Background(), model.DeleteRoomInput{
		RoomID:  roomModel.ID,
		OwnerID: memberUser.ID,
	})
	require.ErrorIs(t, err, roomService.ErrNotOwner)

	persisted, err := subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{
		ID:     roomModel.ID,
		UserID: owner.ID,
	})
	require.NoError(t, err)
	require.Equal(t, owner.ID, persisted.OwnerID)
	require.Equal(t, roomService.DEFAULT_MAX_PLAYERS, persisted.MaxPlayers)
}

func TestRoomServiceOwnerLeaveTransfersOwnershipAndCreatesEvent(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	firstMember := createRoomTestUser(t, subject)
	secondMember := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	time.Sleep(5 * time.Millisecond)
	addRoomTestMember(t, subject, room.ID, firstMember.ID, roomService.ROLE_PLAYER)
	time.Sleep(5 * time.Millisecond)
	addRoomTestMember(t, subject, room.ID, secondMember.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{
		RoomID: room.ID,
		UserID: owner.ID,
	})
	require.NoError(t, err)

	transferredRoom, err := subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{
		ID:     room.ID,
		UserID: firstMember.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstMember.ID, transferredRoom.OwnerID)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: owner.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	events, err := subject.queries.ListRoomEvents(context.Background(), db.ListRoomEventsParams{
		RoomID:     room.ID,
		UserID:     firstMember.ID,
		LimitCount: 10,
	})
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, string(model.EventOwnerTransferred), events[0].EventType)
	require.Equal(t, owner.ID, events[0].ActorID)

	var payload model.OwnerTransferredPayload
	require.NoError(t, json.Unmarshal(events[0].Payload, &payload))
	require.Equal(t, owner.ID, payload.PreviousOwnerID)
	require.Equal(t, firstMember.ID, payload.NewOwnerID)
}

func TestRoomServiceLastMemberLeaveDeletesRoom(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	result, err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{
		RoomID: room.ID,
		UserID: owner.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, result.DeletedRoomID)
	require.Equal(t, room.ID, *result.DeletedRoomID)

	_, err = subject.queries.GetRoomJoinMetaData(context.Background(), room.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	count, err := subject.queries.GetRoomMembersCount(context.Background(), room.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), count)
}

func TestRoomServiceCreatesChatMessagesAndListsEventsOldToNew(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	oldActivity := time.Now().UTC().Add(-2 * time.Hour)
	setRoomLastActivityAt(t, subject, room.ID, oldActivity)

	firstEvent, err := service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: owner.ID,
		Text:    " first message ",
	})
	require.NoError(t, err)
	require.Equal(t, string(model.EventChatMessage), firstEvent.Type)
	require.Equal(t, owner.ID, firstEvent.ActorID)
	var firstPayload model.ChatMessagePayload
	require.NoError(t, json.Unmarshal(firstEvent.Payload, &firstPayload))
	require.Equal(t, "first message", firstPayload.Text)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	time.Sleep(5 * time.Millisecond)
	secondEvent, err := service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: memberUser.ID,
		Text:    "second message",
	})
	require.NoError(t, err)

	events, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: memberUser.ID,
		Limit:  10,
	})
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Equal(t, firstEvent.ID, events[0].ID)
	require.Equal(t, secondEvent.ID, events[1].ID)

	_, err = service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: room.ID,
		UserID: outsider.ID,
		Limit:  10,
	})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	_, err = service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: outsider.ID,
		Text:    "hello",
	})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	_, err = service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: owner.ID,
		Text:    " ",
	})
	require.ErrorIs(t, err, roomService.ErrInvalidInput)
}

func TestRoomServiceChatMessageLengthBoundaries(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	exactLimit := strings.Repeat("x", roomService.MAX_CHAT_MESSAGE_LENGTH)
	event, err := service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: owner.ID,
		Text:    exactLimit,
	})
	require.NoError(t, err)
	var payload model.ChatMessagePayload
	require.NoError(t, json.Unmarshal(event.Payload, &payload))
	require.Len(t, payload.Text, roomService.MAX_CHAT_MESSAGE_LENGTH)

	_, err = service.CreateChatMessage(context.Background(), model.CreateChatMessageInput{
		RoomID:  room.ID,
		ActorID: owner.ID,
		Text:    strings.Repeat("x", roomService.MAX_CHAT_MESSAGE_LENGTH+1),
	})
	require.ErrorIs(t, err, roomService.ErrInvalidInput)
}

func TestRoomServiceListRoomEventsNormalizesLimits(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	payload, err := json.Marshal(model.ChatMessagePayload{Text: "history"})
	require.NoError(t, err)
	for i := 0; i < int(roomService.MAX_ROOM_EVENTS_LIMIT)+5; i++ {
		_, err = subject.queries.CreateRoomEvent(context.Background(), db.CreateRoomEventParams{
			RoomID:    room.ID,
			ActorID:   owner.ID,
			EventType: string(model.EventChatMessage),
			Payload:   payload,
		})
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		limit     int32
		wantCount int
	}{
		{name: "zero uses default", limit: 0, wantCount: int(roomService.DEFAULT_ROOM_EVENTS_LIMIT)},
		{name: "negative uses default", limit: -1, wantCount: int(roomService.DEFAULT_ROOM_EVENTS_LIMIT)},
		{name: "one returns one", limit: 1, wantCount: 1},
		{name: "large caps at max", limit: roomService.MAX_ROOM_EVENTS_LIMIT + 100, wantCount: int(roomService.MAX_ROOM_EVENTS_LIMIT)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
				RoomID: room.ID,
				UserID: owner.ID,
				Limit:  tt.limit,
			})
			require.NoError(t, err)
			require.Len(t, events, tt.wantCount)
		})
	}
}

func TestRoomServiceCreateDiceRollRoomEventMembershipPayloadAndActivity(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.CreateDiceRollRoomEvent(context.Background(), model.CreateDiceRollRoomEventInput{
		RoomID:      room.ID,
		ActorID:     outsider.ID,
		RollID:      "roll-outsider",
		CharacterID: "character-outsider",
		Expression:  "1d20",
		Result:      1,
		Details:     []byte(`[]`),
	})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	oldActivity := time.Now().UTC().Add(-2 * time.Hour)
	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	event, err := service.CreateDiceRollRoomEvent(context.Background(), model.CreateDiceRollRoomEventInput{
		RoomID:      room.ID,
		ActorID:     memberUser.ID,
		RollID:      "roll-member",
		CharacterID: "character-member",
		Expression:  "2d6+1",
		Result:      9,
		Details:     []byte(`[{"type":"dice","sides":6,"rolls":[4,4]},{"type":"modifier","value":1}]`),
	})
	require.NoError(t, err)
	require.Equal(t, string(model.EventDiceRoll), event.Type)
	require.Equal(t, memberUser.ID, event.ActorID)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	var payloadModel model.DiceRollPayload
	require.NoError(t, json.Unmarshal(event.Payload, &payloadModel))
	require.Equal(t, "roll-member", payloadModel.RollID)
	require.Equal(t, "character-member", payloadModel.CharacterID)
	require.Equal(t, "2d6+1", payloadModel.Expression)
	require.Equal(t, int32(9), payloadModel.Result)
	require.JSONEq(t, `[{"type":"dice","sides":6,"rolls":[4,4]},{"type":"modifier","value":1}]`, string(payloadModel.Details))
}

func TestRoomServiceListSelectedCharactersAppliesRoleVisibility(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	gm := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	noCharacterPlayer := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, gm.ID)
	addRoomTestMember(t, subject, room.ID, gm.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, firstPlayer.ID, roomService.ROLE_PLAYER)
	addRoomTestMember(t, subject, room.ID, secondPlayer.ID, roomService.ROLE_PLAYER)
	addRoomTestMember(t, subject, room.ID, noCharacterPlayer.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	gmCharacter := createRoomTestCharacter(t, subject, gm.ID)
	firstPlayerCharacter := createRoomTestCharacter(t, subject, firstPlayer.ID)
	secondPlayerCharacter := createRoomTestCharacter(t, subject, secondPlayer.ID)

	_, err := service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      gm.ID,
		CharacterID: gmCharacter.ID,
	})
	require.NoError(t, err)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      firstPlayer.ID,
		CharacterID: firstPlayerCharacter.ID,
	})
	require.NoError(t, err)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      secondPlayer.ID,
		CharacterID: secondPlayerCharacter.ID,
	})
	require.NoError(t, err)

	updatedFirstPlayerCharacter, err := subject.queries.UpdateCharacter(context.Background(), db.UpdateCharacterParams{
		UserID: firstPlayer.ID,
		ID:     firstPlayerCharacter.ID,
		Name:   "Current DB Investigator",
	})
	require.NoError(t, err)
	currentHP := int16(7)
	maxHP := int16(10)
	_, err = subject.queries.UpsertHealthState(context.Background(), db.UpsertHealthStateParams{
		UserID:      firstPlayer.ID,
		CharacterID: firstPlayerCharacter.ID,
		MaxHp:       &maxHP,
		CurrentHp:   &currentHP,
	})
	require.NoError(t, err)

	playerCharacters, err := service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: firstPlayer.ID,
	})
	require.NoError(t, err)
	require.Len(t, playerCharacters, 1)
	require.Equal(t, firstPlayer.ID, playerCharacters[0].UserID)
	require.Equal(t, roomService.ROLE_PLAYER, playerCharacters[0].Role)
	require.Equal(t, updatedFirstPlayerCharacter.ID, playerCharacters[0].Character.ID)
	require.Equal(t, "Current DB Investigator", playerCharacters[0].Character.Name)
	require.Equal(t, int16(10), playerCharacters[0].Character.HP.MaxHp)
	require.Equal(t, int16(7), playerCharacters[0].Character.HP.CurrentHp)

	gmCharacters, err := service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: gm.ID,
	})
	require.NoError(t, err)
	require.Len(t, gmCharacters, 3)
	requireSelectedCharacterUsers(t, gmCharacters, gm.ID, firstPlayer.ID, secondPlayer.ID)
	require.NotContains(t, selectedCharacterUsers(gmCharacters), noCharacterPlayer.ID)

	_, err = service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: outsider.ID,
	})
	require.ErrorIs(t, err, roomService.ErrNotMember)
}

func TestRoomServiceListSelectedCharactersHidesDeletedCharacterSelection(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	gm := createRoomTestUser(t, subject)
	player := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, gm.ID)
	addRoomTestMember(t, subject, room.ID, gm.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, player.ID, roomService.ROLE_PLAYER)
	character := createRoomTestCharacter(t, subject, player.ID)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      player.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: player.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	member, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: player.ID})
	require.NoError(t, err)
	require.False(t, member.CharacterID.Valid)

	characters, err := service.ListSelectedCharacters(context.Background(), model.ListSelectedCharactersInput{
		RoomID: room.ID,
		UserID: gm.ID,
	})
	require.NoError(t, err)
	require.Empty(t, characters)
}

func TestRoomServiceCreateCharacterChangedRoomEventsPersistsForSelectedRooms(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	character := createRoomTestCharacter(t, subject, owner.ID)
	firstRoom := createRoomTestRoom(t, subject, owner.ID)
	secondRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, firstRoom.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, secondRoom.ID, owner.ID, roomService.ROLE_GM)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      firstRoom.ID,
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      secondRoom.ID,
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	oldActivity := time.Now().UTC().Add(-2 * time.Hour)
	setRoomLastActivityAt(t, subject, firstRoom.ID, oldActivity)
	setRoomLastActivityAt(t, subject, secondRoom.ID, oldActivity)

	sourceEvent := "character.health.upsert_succeeded"
	createdEvents, err := service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		CharacterID: character.ID,
		ActorID:     owner.ID,
		Change: model.CharacterChangedRoomEventChange{
			Resource:    "health",
			Action:      "upsert",
			SourceEvent: &sourceEvent,
		},
	})
	require.NoError(t, err)
	require.Len(t, createdEvents, 2)
	for _, event := range createdEvents {
		require.Equal(t, owner.ID, event.ActorID)
		require.Equal(t, string(model.EventCharacterChanged), event.Type)
	}

	requireRoomCharacterChangedEvent(t, subject, firstRoom.ID, owner.ID, character.ID.String(), "health", "upsert", nil, &sourceEvent)
	requireRoomCharacterChangedEvent(t, subject, secondRoom.ID, owner.ID, character.ID.String(), "health", "upsert", nil, &sourceEvent)
	requireRoomLastActivityAfter(t, subject, firstRoom.ID, owner.ID, oldActivity)
	requireRoomLastActivityAfter(t, subject, secondRoom.ID, owner.ID, oldActivity)
}

func TestRoomServiceCreateCharacterChangedRoomEventsNoOpsForUnselectedCharacter(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	character := createRoomTestCharacter(t, subject, owner.ID)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	sourceEvent := "character.health.upsert_succeeded"
	createdEvents, err := service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		CharacterID: character.ID,
		ActorID:     owner.ID,
		Change: model.CharacterChangedRoomEventChange{
			Resource:    "health",
			Action:      "upsert",
			SourceEvent: &sourceEvent,
		},
	})
	require.NoError(t, err)
	require.Empty(t, createdEvents)

	events, err := subject.queries.ListRoomEvents(context.Background(), db.ListRoomEventsParams{
		RoomID:     room.ID,
		UserID:     owner.ID,
		LimitCount: 10,
	})
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestRoomOwnerCannotEditAnotherUsersCharacterThroughCharacterService(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	player := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, player.ID, roomService.ROLE_PLAYER)
	playerCharacter := createRoomTestCharacter(t, subject, player.ID)
	roomSvc := roomService.NewRoomService(repository.NewRepository(subject.pool))
	_, err := roomSvc.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      player.ID,
		CharacterID: playerCharacter.ID,
	})
	require.NoError(t, err)

	maxHP := int16(10)
	currentHP := int16(7)
	characterSvc := characterService.NewCharacterService(repository.NewRepository(subject.pool))
	_, err = characterSvc.UpsertHealth(context.Background(), healthDTO.UpsertHealthInput{
		UserID:      owner.ID,
		CharacterID: playerCharacter.ID,
		MaxHp:       &maxHP,
		CurrentHp:   &currentHP,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomServiceMutationsTouchRoomActivity(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	targetUser := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)

	oldActivity := time.Now().UTC().Add(-2 * time.Hour)

	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	_, err := service.TransferOwnership(context.Background(), model.TransferOwnershipInput{
		RoomID:     room.ID,
		OwnerID:    owner.ID,
		NewOwnerID: memberUser.ID,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	_, err = service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID:      room.ID,
		UserID:      targetUser.ID,
		InviteToken: room.InviteToken,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	err = service.KickMember(context.Background(), model.KickMemberInput{
		RoomID:       room.ID,
		ActorUserID:  memberUser.ID,
		TargetUserID: targetUser.ID,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	_, err = service.ChangeRole(context.Background(), model.ChangeRoleInput{
		RoomID:       room.ID,
		ActorUserID:  memberUser.ID,
		TargetUserID: owner.ID,
		Role:         roomService.ROLE_PLAYER,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	character := createRoomTestCharacter(t, subject, owner.ID)
	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, owner.ID, oldActivity)

	setRoomLastActivityAt(t, subject, room.ID, oldActivity)
	_, err = service.LeaveRoom(context.Background(), model.LeaveRoomInput{
		RoomID: room.ID,
		UserID: owner.ID,
	})
	require.NoError(t, err)
	requireRoomLastActivityAfter(t, subject, room.ID, memberUser.ID, oldActivity)
}

func TestRoomServiceCleanupRoomsDeletesInactiveAndInvalidRooms(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	inactiveRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, inactiveRoom.ID, owner.ID, roomService.ROLE_GM)
	setRoomLastActivityAt(t, subject, inactiveRoom.ID, time.Now().UTC().Add(-13*time.Hour))

	activeRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, activeRoom.ID, owner.ID, roomService.ROLE_GM)
	setRoomLastActivityAt(t, subject, activeRoom.ID, time.Now().UTC())

	noMembersRoom := createRoomTestRoom(t, subject, owner.ID)
	setRoomLastActivityAt(t, subject, noMembersRoom.ID, time.Now().UTC())

	ownerNotMemberRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, ownerNotMemberRoom.ID, memberUser.ID, roomService.ROLE_PLAYER)
	setRoomLastActivityAt(t, subject, ownerNotMemberRoom.ID, time.Now().UTC())

	result, err := service.CleanupRooms(context.Background(), model.CleanupRoomsInput{Now: time.Now().UTC()})
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.InactiveDeleted, 1)
	require.GreaterOrEqual(t, result.InvalidDeleted, 2)
	requireCleanupDeletedRoomIDs(t, result, inactiveRoom.ID, noMembersRoom.ID, ownerNotMemberRoom.ID)

	_, err = subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: inactiveRoom.ID, UserID: owner.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetRoomJoinMetaData(context.Background(), noMembersRoom.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetRoomJoinMetaData(context.Background(), ownerNotMemberRoom.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: activeRoom.ID, UserID: owner.ID})
	require.NoError(t, err)
}

func TestRoomServiceUpdateRoomValidatesMaxPlayersAgainstCurrentMembers(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:     room.ID,
		OwnerID:    owner.ID,
		MaxPlayers: 1,
	})
	require.ErrorIs(t, err, roomService.ErrInvalidInput)

	updated, err := service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:     room.ID,
		OwnerID:    owner.ID,
		MaxPlayers: 2,
	})
	require.NoError(t, err)
	require.Equal(t, int32(2), updated.MaxPlayers)
}

func TestRoomServiceUpdateRoomCanChangePassword(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	player := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:  owner.ID,
		Password: "old-password",
	})
	require.NoError(t, err)

	blankPassword := "   "
	_, err = service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:     roomModel.ID,
		OwnerID:    owner.ID,
		MaxPlayers: roomModel.MaxPlayers,
		Password:   &blankPassword,
	})
	require.ErrorIs(t, err, roomService.ErrInvalidPassword)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: roomModel.ID, Password: "old-password", UserID: player.ID},
	)
	require.NoError(t, err)

	newPassword := "new-password"
	_, err = service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:     roomModel.ID,
		OwnerID:    owner.ID,
		MaxPlayers: roomModel.MaxPlayers,
		Password:   &newPassword,
	})
	require.NoError(t, err)

	_, err = service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: roomModel.ID, Password: "old-password", UserID: player.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomNotFound)

	memberModel, err := service.JoinRoom(
		context.Background(),
		model.JoinRoomInput{RoomID: roomModel.ID, Password: "new-password", UserID: secondPlayer.ID},
	)
	require.NoError(t, err)
	require.Equal(t, secondPlayer.ID, memberModel.UserID)
}

func TestRoomServiceMapsNoRowsForMembershipOperations(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	err = service.KickMember(
		context.Background(),
		model.KickMemberInput{RoomID: room.ID, ActorUserID: memberUser.ID, TargetUserID: owner.ID},
	)
	require.ErrorIs(t, err, roomService.ErrNotMember)

	err = service.KickMember(
		context.Background(),
		model.KickMemberInput{RoomID: room.ID, ActorUserID: owner.ID, TargetUserID: owner.ID},
	)
	require.ErrorIs(t, err, roomService.ErrCannotKickOwner)

	_, err = service.ChangeRole(
		context.Background(),
		model.ChangeRoleInput{RoomID: room.ID, ActorUserID: owner.ID, TargetUserID: memberUser.ID, Role: roomService.ROLE_GM},
	)
	require.ErrorIs(t, err, roomService.ErrNotMember)

	_, err = service.ChangeRole(
		context.Background(),
		model.ChangeRoleInput{RoomID: room.ID, ActorUserID: owner.ID, TargetUserID: memberUser.ID, Role: "keeper"},
	)
	require.ErrorIs(t, err, roomService.ErrInvalidInput)

	_, err = service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: roomTestUUID("33333333-3333-3333-3333-333333333333"),
	})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomServiceSelectCharacterRejectsCharacterNotOwnedForMember(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	otherUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)
	otherCharacter := createRoomTestCharacter(t, subject, otherUser.ID)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err := service.SelectCharacter(context.Background(), model.SelectCharacterInput{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: otherCharacter.ID,
	})
	require.ErrorIs(t, err, roomService.ErrCharacterNotOwned)
}

func TestRoomServiceReusablePermissionChecks(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.ROLE_GM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.ROLE_PLAYER)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	require.NoError(t, service.EnsureMember(context.Background(), room.ID, memberUser.ID))
	require.NoError(t, service.EnsureCanPublishRoomEvent(context.Background(), room.ID, memberUser.ID))
	require.ErrorIs(t, service.EnsureMember(context.Background(), room.ID, outsider.ID), roomService.ErrNotMember)
	require.NoError(t, service.EnsureOwner(context.Background(), room.ID, owner.ID))
	require.ErrorIs(t, service.EnsureOwner(context.Background(), room.ID, memberUser.ID), roomService.ErrNotOwner)
	require.ErrorIs(t, service.EnsureOwner(context.Background(), room.ID, outsider.ID), roomService.ErrNotMember)
}
