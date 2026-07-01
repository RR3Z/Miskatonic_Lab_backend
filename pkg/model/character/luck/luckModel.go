package luck

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type LuckModel struct {
	ID           pgtype.UUID        `json:"id"`
	CharacterID  pgtype.UUID        `json:"character_id"`
	StartingLuck int16              `json:"starting_luck"`
	CurrentLuck  int16              `json:"current_luck"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

func ToLuckModel(l db.LuckState) LuckModel {
	return LuckModel{
		ID:           l.ID,
		CharacterID:  l.CharacterID,
		StartingLuck: l.StartingLuck,
		CurrentLuck:  l.CurrentLuck,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}
