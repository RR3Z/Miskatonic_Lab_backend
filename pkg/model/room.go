package model

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomMemberModel struct {
	ID          pgtype.UUID        `json:"id"`
	RoomID      pgtype.UUID        `json:"room_id"`
	UserID      string             `json:"user_id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Role        string             `json:"role"`
	JoinedAt    pgtype.Timestamptz `json:"joined_at"`
}

type RoomModel struct {
	ID          pgtype.UUID        `json:"id"`
	OwnerID     string             `json:"owner_id"`
	MaxPlayers  int32              `json:"max_players"`
	InviteToken string             `json:"invite_token"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	Members     []RoomMemberModel  `json:"members"`
}

func ToRoomModel(r db.Room, members []db.RoomMember) RoomModel {
	m := make([]RoomMemberModel, len(members))
	for i, mb := range members {
		m[i] = ToRoomMemberModel(mb)
	}

	return RoomModel{
		ID:          r.ID,
		OwnerID:     r.OwnerID,
		MaxPlayers:  r.MaxPlayers,
		InviteToken: r.InviteToken,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Members:     m,
	}
}

func ToRoomMemberModel(m db.RoomMember) RoomMemberModel {
	return RoomMemberModel{
		ID:          m.ID,
		RoomID:      m.RoomID,
		UserID:      m.UserID,
		CharacterID: m.CharacterID,
		Role:        m.Role,
		JoinedAt:    m.JoinedAt,
	}
}

type CreateRoomRequest struct {
	MaxPlayers *int32 `json:"max_players"`
}

type UpdateRoomRequest struct {
	MaxPlayers int32 `json:"max_players"`
}

type SelectCharacterRequest struct {
	CharacterID pgtype.UUID `json:"character_id"`
}

type JoinRoomRequest struct {
	InviteToken string `json:"invite_token"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type TransferRoomOwnershipRequest struct {
	UserID string `json:"user_id"`
}
