package magicDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type MagicModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	MaxMp       int16              `json:"max_mp"`
	CurrentMp   int16              `json:"current_mp"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func ToMagicModel(m db.MagicState) MagicModel {
	return MagicModel{
		ID:          m.ID,
		CharacterID: m.CharacterID,
		MaxMp:       m.MaxMp,
		CurrentMp:   m.CurrentMp,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
