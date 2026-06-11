package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestRoomMemberCrudRoleAndConstraints(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	invalidRoleUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)

	ownerMember := addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	require.Equal(t, "gm", ownerMember.Role)
	time.Sleep(5 * time.Millisecond)
	member := addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")
	require.Equal(t, "player", member.Role)

	_, err := subject.queries.AddMember(context.Background(), db.AddMemberParams{RoomID: room.ID, UserID: memberUser.ID, Role: "player"})
	requireRoomPostgresErrorCode(t, err, "23505")

	_, err = subject.queries.AddMember(context.Background(), db.AddMemberParams{RoomID: roomTestUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), UserID: memberUser.ID, Role: "player"})
	requireRoomPostgresErrorCode(t, err, "23503")

	_, err = subject.queries.AddMember(context.Background(), db.AddMemberParams{RoomID: room.ID, UserID: "missing_user", Role: "player"})
	requireRoomPostgresErrorCode(t, err, "23503")

	_, err = subject.queries.AddMember(context.Background(), db.AddMemberParams{RoomID: room.ID, UserID: invalidRoleUser.ID, Role: "keeper"})
	requireRoomPostgresErrorCode(t, err, "23514")

	fetched, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, member.ID, fetched.ID)

	members, err := subject.queries.ListMembersByRoomID(context.Background(), db.ListMembersByRoomIDParams{RoomID: room.ID, UserID: owner.ID})
	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, ownerMember.ID, members[0].ID)
	require.Equal(t, member.ID, members[1].ID)

	updated, err := subject.queries.UpdateMemberRole(context.Background(), db.UpdateMemberRoleParams{RoomID: room.ID, UserID: memberUser.ID, Role: "gm"})
	require.NoError(t, err)
	require.Equal(t, "gm", updated.Role)

	_, err = subject.queries.UpdateMemberRole(context.Background(), db.UpdateMemberRoleParams{RoomID: room.ID, UserID: memberUser.ID, Role: "keeper"})
	requireRoomPostgresErrorCode(t, err, "23514")
	fetched, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, "gm", fetched.Role)

	removed, err := subject.queries.RemoveMember(context.Background(), db.RemoveMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, member.ID, removed.ID)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.RemoveMember(context.Background(), db.RemoveMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomMemberSelectCharacterOwnershipAndCharacterDelete(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	otherUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")
	character := createRoomTestCharacter(t, subject, memberUser.ID)
	otherCharacter := createRoomTestCharacter(t, subject, otherUser.ID)

	updated, err := subject.queries.UpdateMemberCharacter(context.Background(), db.UpdateMemberCharacterParams{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, character.ID, updated.CharacterID)

	_, err = subject.queries.UpdateMemberCharacter(context.Background(), db.UpdateMemberCharacterParams{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: otherCharacter.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetched, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, character.ID, fetched.CharacterID)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{UserID: memberUser.ID, ID: character.ID})
	require.NoError(t, err)

	fetched, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.False(t, fetched.CharacterID.Valid)
}

func TestRoomMemberUserDeleteRemovesMembership(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")

	err := subject.queries.DeleteUserByClerkID(context.Background(), memberUser.ID)
	require.NoError(t, err)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomOwnershipTransferRequiresOwnerAndMemberWithoutChangingRole(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	outsider := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	member := addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")

	transferred, err := subject.queries.TransferRoomOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    owner.ID,
		NewOwnerID: memberUser.ID,
	})
	require.NoError(t, err)
	require.Equal(t, memberUser.ID, transferred.OwnerID)

	fetchedMember, err := subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, member.Role, fetchedMember.Role)

	_, err = subject.queries.TransferRoomOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    memberUser.ID,
		NewOwnerID: outsider.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.TransferRoomOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    outsider.ID,
		NewOwnerID: owner.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	selfTransfer, err := subject.queries.TransferRoomOwnership(context.Background(), db.TransferRoomOwnershipParams{
		ID:         room.ID,
		OwnerID:    memberUser.ID,
		NewOwnerID: memberUser.ID,
	})
	require.NoError(t, err)
	require.Equal(t, memberUser.ID, selfTransfer.OwnerID)

	fetchedMember, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.NoError(t, err)
	require.Equal(t, "player", fetchedMember.Role)
}

func TestRoomDeleteCascadesMembers(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")

	_, err := subject.queries.DeleteRoom(context.Background(), db.DeleteRoomParams{ID: room.ID, OwnerID: owner.ID})
	require.NoError(t, err)

	_, err = subject.queries.GetMember(context.Background(), db.GetMemberParams{RoomID: room.ID, UserID: memberUser.ID})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestRoomMemberMissingCharacterUpdateReturnsNoRows(t *testing.T) {
	subject := newRoomIntegrationSubject(t)
	owner := createRoomTestUser(t, subject)
	memberUser := createRoomTestUser(t, subject)
	room := createRoomTestRoom(t, subject, owner.ID)
	addRoomTestMember(t, subject, room.ID, owner.ID, "gm")
	addRoomTestMember(t, subject, room.ID, memberUser.ID, "player")

	_, err := subject.queries.UpdateMemberCharacter(context.Background(), db.UpdateMemberCharacterParams{
		RoomID:      room.ID,
		UserID:      memberUser.ID,
		CharacterID: pgtype.UUID{},
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
