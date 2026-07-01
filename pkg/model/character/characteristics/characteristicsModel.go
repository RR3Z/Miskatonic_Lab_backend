package characteristicsDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacteristicsModel struct {
	ID           pgtype.UUID        `json:"id"`
	CharacterID  pgtype.UUID        `json:"character_id"`
	Strength     *int16             `json:"strength"`
	Constitution *int16             `json:"constitution"`
	Size         *int16             `json:"size"`
	Dexterity    *int16             `json:"dexterity"`
	Appearance   *int16             `json:"appearance"`
	Intelligence *int16             `json:"intelligence"`
	Power        *int16             `json:"power"`
	Education    *int16             `json:"education"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

func ToCharacteristicsModel(c db.Characteristic) CharacteristicsModel {
	return CharacteristicsModel{
		ID:           c.ID,
		CharacterID:  c.CharacterID,
		Strength:     c.Strength,
		Constitution: c.Constitution,
		Size:         c.Size,
		Dexterity:    c.Dexterity,
		Appearance:   c.Appearance,
		Intelligence: c.Intelligence,
		Power:        c.Power,
		Education:    c.Education,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
