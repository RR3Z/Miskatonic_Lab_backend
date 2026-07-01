package luck

import "github.com/jackc/pgx/v5/pgtype"

type LuckModel struct {
	ID           pgtype.UUID        `json:"id"`
	CharacterID  pgtype.UUID        `json:"character_id"`
	StartingLuck int16              `json:"starting_luck"`
	CurrentLuck  int16              `json:"current_luck"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}
