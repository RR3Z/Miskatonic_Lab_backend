package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestRoomServiceCreateRoomCreatesOwnerMemberAndInviteToken(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	roomModel, err := service.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:    owner.ID,
		MaxPlayers: 3,
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
}

func TestRoomServiceJoinRoomRequiresInviteAndCapacity(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	firstPlayer := createRoomTestUser(t, subject)
	secondPlayer := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	room, err := subject.queries.UpdateRoom(context.Background(), db.UpdateRoomParams{ID: room.ID, OwnerID: owner.ID, MaxPlayers: 2})
	require.NoError(t, err)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	_, err = service.JoinRoom(
		context.Background(),
		db.GetRoomMetaDataParams{ID: room.ID, InviteToken: "wrong_token"},
		db.GetMemberParams{RoomID: room.ID, UserID: firstPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomNotFound)

	memberModel, err := service.JoinRoom(
		context.Background(),
		db.GetRoomMetaDataParams{ID: room.ID, InviteToken: room.InviteToken},
		db.GetMemberParams{RoomID: room.ID, UserID: firstPlayer.ID},
	)
	require.NoError(t, err)
	require.Equal(t, firstPlayer.ID, memberModel.UserID)
	require.Equal(t, "player", memberModel.Role)

	_, err = service.JoinRoom(
		context.Background(),
		db.GetRoomMetaDataParams{ID: room.ID, InviteToken: room.InviteToken},
		db.GetMemberParams{RoomID: room.ID, UserID: firstPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrAlreadyMember)

	_, err = service.JoinRoom(
		context.Background(),
		db.GetRoomMetaDataParams{ID: room.ID, InviteToken: room.InviteToken},
		db.GetMemberParams{RoomID: room.ID, UserID: secondPlayer.ID},
	)
	require.ErrorIs(t, err, roomService.ErrRoomFull)
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

	transferred, err := service.TransferOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    owner.ID,
		NewOwnerID: memberUser.ID,
	})
	require.NoError(t, err)
	require.Equal(t, memberUser.ID, transferred.OwnerID)

	member, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, "player", member.Role)

	_, err = service.TransferOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    memberUser.ID,
		NewOwnerID: outsider.ID,
	})
	require.ErrorIs(t, err, roomService.ErrNotOwner)
}

func TestRoomServiceMapsNoRowsForMembershipOperations(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	service := roomService.NewRoomService(repository.NewRepository(subject.pool))

	err := service.LeaveRoom(context.Background(), db.RemoveMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, roomService.ErrNotMember)

	err = service.KickMember(
		context.Background(),
		db.GetRoomByIDParams{ID: room.ID, UserID: memberUser.ID},
		db.RemoveMemberParams{RoomID: room.ID, UserID: owner.ID},
	)
	require.ErrorIs(t, err, roomService.ErrNotMember)

	err = service.KickMember(
		context.Background(),
		db.GetRoomByIDParams{ID: room.ID, UserID: owner.ID},
		db.RemoveMemberParams{RoomID: room.ID, UserID: owner.ID},
	)
	require.ErrorIs(t, err, roomService.ErrCannotKickOwner)

	_, err = service.ChangeRole(
		context.Background(),
		db.GetRoomByIDParams{ID: room.ID, UserID: owner.ID},
		db.UpdateMemberRoleParams{RoomID: room.ID, UserID: memberUser.ID, Role: "gm"},
	)
	require.ErrorIs(t, err, roomService.ErrNotMember)

	_, err = service.SelectCharacter(context.Background(), db.UpdateMemberCharacterParams{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: roomTestUUID("33333333-3333-3333-3333-333333333333"),
	})
	require.ErrorIs(t, err, roomService.ErrCharacterNotOwned)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
