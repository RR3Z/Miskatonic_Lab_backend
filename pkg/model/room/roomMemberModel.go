package roomDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomMemberModel struct {
	ID          pgtype.UUID        `json:"id"`
	RoomID      pgtype.UUID        `json:"room_id"`
	UserID      string             `json:"user_id"`
	Username    string             `json:"username,omitempty"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Role        string             `json:"role"`
	JoinedAt    pgtype.Timestamptz `json:"joined_at"`
}

func ToRoomMemberModels(members []db.ListMembersByRoomIDRow) []RoomMemberModel {
	models := make([]RoomMemberModel, len(members))
	for index, member := range members {
		models[index] = RoomMemberModel{
			ID:          member.ID,
			RoomID:      member.RoomID,
			UserID:      member.UserID,
			Username:    member.Username,
			CharacterID: member.CharacterID,
			Role:        member.Role,
			JoinedAt:    member.JoinedAt,
		}
	}

	return models
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
