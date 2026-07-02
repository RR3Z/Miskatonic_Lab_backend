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
	require.Equal(t, roomService.DefaultMaxPlayers, roomModel.MaxPlayers)

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
	require.Equal(t, roomService.RolePlayer, memberModel.Role)
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
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
	time.Sleep(5 * time.Millisecond)
	addRoomTestMember(t, subject, room.ID, firstMember.ID, roomService.RolePlayer)
	time.Sleep(5 * time.Millisecond)
	addRoomTestMember(t, subject, room.ID, secondMember.ID, roomService.RolePlayer)
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
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
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

func TestRoomServiceMutationsTouchRoomActivity(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	targetUser := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.RolePlayer)

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
		Role:         roomService.RolePlayer,
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
	addRoomTestMember(t, subject, inactiveRoom.ID, owner.ID, roomService.RoleGM)
	setRoomLastActivityAt(t, subject, inactiveRoom.ID, time.Now().UTC().Add(-13*time.Hour))

	activeRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, activeRoom.ID, owner.ID, roomService.RoleGM)
	setRoomLastActivityAt(t, subject, activeRoom.ID, time.Now().UTC())

	noMembersRoom := createRoomTestRoom(t, subject, owner.ID)
	setRoomLastActivityAt(t, subject, noMembersRoom.ID, time.Now().UTC())

	ownerNotMemberRoom := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, ownerNotMemberRoom.ID, memberUser.ID, roomService.RolePlayer)
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
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.RolePlayer)
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
		model.ChangeRoleInput{RoomID: room.ID, ActorUserID: owner.ID, TargetUserID: memberUser.ID, Role: roomService.RoleGM},
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
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.RolePlayer)
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
	addRoomTestMember(t, subject, room.ID, owner.ID, roomService.RoleGM)
	addRoomTestMember(t, subject, room.ID, memberUser.ID, roomService.RolePlayer)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	require.NoError(t, service.EnsureMember(context.Background(), room.ID, memberUser.ID))
	require.NoError(t, service.EnsureCanPublishRoomEvent(context.Background(), room.ID, memberUser.ID))
	require.ErrorIs(t, service.EnsureMember(context.Background(), room.ID, outsider.ID), roomService.ErrNotMember)
	require.NoError(t, service.EnsureOwner(context.Background(), room.ID, owner.ID))
	require.ErrorIs(t, service.EnsureOwner(context.Background(), room.ID, memberUser.ID), roomService.ErrNotOwner)
	require.ErrorIs(t, service.EnsureOwner(context.Background(), room.ID, outsider.ID), roomService.ErrNotMember)
}
