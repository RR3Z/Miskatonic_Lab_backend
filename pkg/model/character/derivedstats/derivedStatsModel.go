package derivedstats

import "github.com/jackc/pgx/v5/pgtype"

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
