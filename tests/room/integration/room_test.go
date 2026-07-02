package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestRoomCreateGetUpdateDeleteAndMemberList(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)

	room, err := subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      owner.ID,
		MaxPlayers:   3,
		InviteToken:  "invite_" + uniqueRoomIntegrationSuffix(),
		PasswordHash: "test_password_hash",
	})
	require.NoError(t, err)
	require.True(t, room.ID.Valid)
	require.Equal(t, owner.ID, room.OwnerID)
	require.Equal(t, int32(3), room.MaxPlayers)
	require.NotEmpty(t, room.InviteToken)

	ownerMember := addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	time.Sleep(5 * time.Millisecond)
	member := addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")

	fetched, err := subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: room.ID, UserID: owner.ID})
	require.NoError(t, err)
	require.Equal(t, room.ID, fetched.ID)

	_, err = subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: room.ID, UserID: outsider.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	members, err := subject.queries.ListMembersByRoomID(context.Background(), db.ListMembersByRoomIDParams{RoomID: room.ID, UserID: owner.ID})
	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, ownerMember.ID, members[0].ID)
	require.Equal(t, member.ID, members[1].ID)

	outsiderMembers, err := subject.queries.ListMembersByRoomID(context.Background(), db.ListMembersByRoomIDParams{RoomID: room.ID, UserID: outsider.ID})
	require.NoError(t, err)
	require.Empty(t, outsiderMembers)

	updated, err := subject.queries.UpdateRoom(context.Background(), db.UpdateRoomParams{ID: room.ID, OwnerID: owner.ID, MaxPlayers: 5})
	require.NoError(t, err)
	require.Equal(t, int32(5), updated.MaxPlayers)
	require.Equal(t, room.PasswordHash, updated.PasswordHash)
	require.True(t, updated.UpdatedAt.Time.After(room.UpdatedAt.Time) || updated.UpdatedAt.Time.Equal(room.UpdatedAt.Time))

	_, err = subject.queries.UpdateRoom(context.Background(), db.UpdateRoomParams{ID: room.ID, OwnerID: memberUser.ID, MaxPlayers: 6})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteRoom(context.Background(), db.DeleteRoomParams{ID: room.ID, OwnerID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deleted, err := subject.queries.DeleteRoom(context.Background(), db.DeleteRoomParams{ID: room.ID, OwnerID: owner.ID})
	require.NoError(t, err)
	require.Equal(t, room.ID, deleted.ID)

	_, err = subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: room.ID, UserID: owner.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomConstraintsAndInviteMetadata(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	inviteToken := "invite_" + uniqueRoomIntegrationSuffix()

	room, err := subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      owner.ID,
		MaxPlayers:   2,
		InviteToken:  inviteToken,
		PasswordHash: "test_password_hash",
	})
	require.NoError(t, err)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")

	_, err = subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      "missing_user",
		MaxPlayers:   2,
		InviteToken:  "invite_" + uniqueRoomIntegrationSuffix(),
		PasswordHash: "test_password_hash",
	})
	requireRoomPostgresErrorCode(t, err, "23503")

	_, err = subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      owner.ID,
		MaxPlayers:   0,
		InviteToken:  "invite_" + uniqueRoomIntegrationSuffix(),
		PasswordHash: "test_password_hash",
	})
	requireRoomPostgresErrorCode(t, err, "23514")

	_, err = subject.queries.CreateRoom(context.Background(), db.CreateRoomParams{
		OwnerID:      owner.ID,
		MaxPlayers:   2,
		InviteToken:  inviteToken,
		PasswordHash: "test_password_hash",
	})
	requireRoomPostgresErrorCode(t, err, "23505")

	meta, err := subject.queries.GetRoomMetaData(context.Background(), db.GetRoomMetaDataParams{ID: room.ID, InviteToken: inviteToken})
	require.NoError(t, err)
	require.Equal(t, room.ID, meta.ID)
	require.Equal(t, int32(2), meta.MaxPlayers)

	joinMeta, err := subject.queries.GetRoomJoinMetaData(context.Background(), room.ID)
	require.NoError(t, err)
	require.Equal(t, room.ID, joinMeta.ID)
	require.Equal(t, inviteToken, joinMeta.InviteToken)
	require.Equal(t, "test_password_hash", joinMeta.PasswordHash)

	_, err = subject.queries.GetRoomMetaData(context.Background(), db.GetRoomMetaDataParams{ID: room.ID, InviteToken: "wrong_token"})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	count, err := subject.queries.GetRoomMembersCount(context.Background(), room.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), count)
}

func TestRoomUserDeleteCascadesOwnedRoomsAndMemberships(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	member := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, member.ID, "player")

	err := subject.queries.DeleteUserByClerkID(context.Background(), owner.ID)
	require.NoError(t, err)

	count, err := subject.queries.GetRoomMembersCount(context.Background(), room.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), count)
	_, err = subject.queries.GetRoomByID(context.Background(), db.GetRoomByIDParams{ID: room.ID, UserID: member.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
