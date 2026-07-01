package magic

import "github.com/jackc/pgx/v5/pgtype"

type MagicModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	MaxMp       int16              `json:"max_mp"`
	CurrentMp   int16              `json:"current_mp"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
