package sanityDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SanityModel struct {
	ID            pgtype.UUID        `json:"id"`
	CharacterID   pgtype.UUID        `json:"character_id"`
	MaxSanity     int16              `json:"max_sanity"`
	CurrentSanity int16              `json:"current_sanity"`
	TempInsanity  bool               `json:"temp_insanity"`
	IndefInsanity bool               `json:"indef_insanity"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
}

func ToSanityModel(s db.SanityState) SanityModel {
	return SanityModel{
		ID:            s.ID,
		CharacterID:   s.CharacterID,
		MaxSanity:     s.MaxSanity,
		CurrentSanity: s.CurrentSanity,
		TempInsanity:  s.TempInsanity,
		IndefInsanity: s.IndefInsanity,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
