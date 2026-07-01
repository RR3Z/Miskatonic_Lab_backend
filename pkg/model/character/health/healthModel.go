package health

import "github.com/jackc/pgx/v5/pgtype"

type HealthModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	MaxHp       int16              `json:"max_hp"`
	CurrentHp   int16              `json:"current_hp"`
	MajorWound  bool               `json:"major_wound"`
	Unconscious bool               `json:"unconscious"`
	Dying       bool               `json:"dying"`
	Dead        bool               `json:"dead"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
