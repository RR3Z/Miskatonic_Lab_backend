package roomDTO

import (
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SelectedCharacterModel struct {
	MemberID  pgtype.UUID                 `json:"member_id"`
	UserID    string                      `json:"user_id"`
	Role      string                      `json:"role"`
	Character characterDTO.CharacterModel `json:"character"`
}

func ToSelectedCharacterModel(member db.RoomMember, character characterDTO.CharacterModel) SelectedCharacterModel {
	return SelectedCharacterModel{
		MemberID:  member.ID,
		UserID:    member.UserID,
		Role:      member.Role,
		Character: character,
	}
}
