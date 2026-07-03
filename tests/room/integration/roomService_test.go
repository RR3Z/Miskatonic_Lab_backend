package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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

	err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{
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
	require.Equal(t, string(roomEvents.EventOwnerTransferred), events[0].EventType)
	require.Equal(t, owner.ID, events[0].ActorID)

	var payload roomEvents.OwnerTransferredPayload
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

	err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{
		RoomID: room.ID,
		UserID: owner.ID,
	})
	require.NoError(t, err)

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
	require.Equal(t, string(roomEvents.EventChatMessage), firstEvent.Type)
	require.Equal(t, owner.ID, firstEvent.ActorID)
	var firstPayload roomEvents.ChatMessagePayload
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

func requireSelectedCharacterUsers(t *testing.T, characters []model.SelectedCharacterModel, expectedUsers ...string) {
	t.Helper()

	users := selectedCharacterUsers(characters)
	for _, userID := range expectedUsers {
		require.Contains(t, users, userID)
	}
}

func selectedCharacterUsers(characters []model.SelectedCharacterModel) []string {
	users := make([]string, 0, len(characters))
	for _, character := range characters {
		users = append(users, character.UserID)
	}
	return users
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
	err = service.LeaveRoom(context.Background(), model.LeaveRoomInput{
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
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), model.CreateRoomInput{
		OwnerID:  owner.ID,
		Password: "old-password",
	})
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
		model.JoinRoomInput{RoomID: roomModel.ID, Password: "new-password", UserID: player.ID},
	)
	require.NoError(t, err)
	require.Equal(t, player.ID, memberModel.UserID)
}

func TestRoomServiceMapsNoRowsForMembershipOperations(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{RoomID: room.ID, UserID: memberUser.ID})
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
