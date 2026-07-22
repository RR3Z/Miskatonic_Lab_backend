package tests

import (
	"testing"
	"time"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestToRoomModelMapsRoomAndMembers(t *testing.T) {
	roomID := testUUID("11111111-1111-1111-1111-111111111111")
	memberID := testUUID("22222222-2222-2222-2222-222222222222")
	characterID := testUUID("33333333-3333-3333-3333-333333333333")
	createdAt := pgtype.Timestamptz{Time: time.Date(2026, 6, 11, 10, 0, 0, 0, time.UTC), Valid: true}
	updatedAt := pgtype.Timestamptz{Time: time.Date(2026, 6, 11, 11, 0, 0, 0, time.UTC), Valid: true}
	joinedAt := pgtype.Timestamptz{Time: time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC), Valid: true}

	roomModel := model.ToRoomModel(db.Room{
		ID:          roomID,
		OwnerID:     "owner_1",
		MaxPlayers:  5,
		InviteToken: "invite_token",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, []db.RoomMember{{
		ID:          memberID,
		RoomID:      roomID,
		UserID:      "member_1",
		CharacterID: characterID,
		Role:        "player",
		JoinedAt:    joinedAt,
	}}, "owner_1")

	require.Equal(t, roomID, roomModel.ID)
	require.Equal(t, "owner_1", roomModel.OwnerID)
	require.Equal(t, int32(5), roomModel.MaxPlayers)
	require.Equal(t, "invite_token", roomModel.InviteToken)
	require.Equal(t, createdAt, roomModel.CreatedAt)
	require.Equal(t, updatedAt, roomModel.UpdatedAt)
	require.Len(t, roomModel.Members, 1)
	require.Equal(t, memberID, roomModel.Members[0].ID)
	require.Equal(t, roomID, roomModel.Members[0].RoomID)
	require.Equal(t, "member_1", roomModel.Members[0].UserID)
	require.Equal(t, characterID, roomModel.Members[0].CharacterID)
	require.Equal(t, "player", roomModel.Members[0].Role)
	require.Equal(t, joinedAt, roomModel.Members[0].JoinedAt)
}

func TestToRoomModelMapsNilMembersToEmptySlice(t *testing.T) {
	roomModel := model.ToRoomModel(db.Room{OwnerID: "owner_1"}, nil, "owner_1")

	require.NotNil(t, roomModel.Members)
	require.Empty(t, roomModel.Members)
}

func TestToRoomMemberModelPreservesInvalidCharacterID(t *testing.T) {
	memberID := testUUID("22222222-2222-2222-2222-222222222222")
	roomID := testUUID("11111111-1111-1111-1111-111111111111")

	memberModel := model.ToRoomMemberModel(db.RoomMember{
		ID:          memberID,
		RoomID:      roomID,
		UserID:      "member_1",
		CharacterID: pgtype.UUID{},
		Role:        "gm",
	})

	require.Equal(t, memberID, memberModel.ID)
	require.Equal(t, roomID, memberModel.RoomID)
	require.Equal(t, "member_1", memberModel.UserID)
	require.False(t, memberModel.CharacterID.Valid)
	require.Equal(t, "gm", memberModel.Role)
}

func TestToRoomEventModelMapsSequence(t *testing.T) {
	roomID := testUUID("11111111-1111-1111-1111-111111111111")
	eventID := testUUID("22222222-2222-2222-2222-222222222222")

	eventModel := model.ToRoomEventModel(db.RoomEvent{
		ID:        eventID,
		RoomID:    roomID,
		Sequence:  17,
		ActorID:   "user_1",
		EventType: string(model.EventChatMessage),
		Payload:   []byte(`{"text":"hello"}`),
	})

	require.Equal(t, int64(17), eventModel.Sequence)
	require.Equal(t, eventID, eventModel.ID)
	require.Equal(t, roomID, eventModel.RoomID)
}

func TestToRoomModelWithUsernamesMapsMemberUsername(t *testing.T) {
	roomID := testUUID("11111111-1111-1111-1111-111111111111")
	memberID := testUUID("22222222-2222-2222-2222-222222222222")

	roomModel := model.ToRoomModelWithUsernames(db.Room{ID: roomID}, []db.ListMembersByRoomIDRow{{
		ID:       memberID,
		RoomID:   roomID,
		Role:     "player",
		UserID:   "user_123",
		Username: "Роберт",
	}}, "user_123")

	require.Len(t, roomModel.Members, 1)
	require.Equal(t, "Роберт", roomModel.Members[0].Username)
}
