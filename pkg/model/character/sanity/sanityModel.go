package sanity

import "github.com/jackc/pgx/v5/pgtype"

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
