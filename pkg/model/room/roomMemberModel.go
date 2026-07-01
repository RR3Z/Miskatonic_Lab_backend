package roomDTO

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
