package derivedStatsDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type DerivedStatsModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Speed       *int16             `json:"speed"`
	Physique    *int16             `json:"physique"`
	DamageBonus *string            `json:"damage_bonus"`
	DodgeValue  *int16             `json:"dodge_value"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func ToDerivedStatsModel(d db.DerivedStat) DerivedStatsModel {
	return DerivedStatsModel{
		ID:          d.ID,
		CharacterID: d.CharacterID,
		Speed:       d.Speed,
		Physique:    d.Physique,
		DamageBonus: d.DamageBonus,
		DodgeValue:  d.DodgeValue,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
