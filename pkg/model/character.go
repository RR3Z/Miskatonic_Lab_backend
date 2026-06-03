package model

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterModel struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`

	Name            string            `json:"name"`
	PlayerName      *string           `json:"player_name"`
	Occupation      *string           `json:"occupation"`
	Age             *int16            `json:"age"`
	Sex             *string           `json:"sex"`
	Residence       *string           `json:"residence"`
	Birthplace      *string           `json:"birthplace"`
	Skills          []SkillModel      `json:"skills"`
	Characteristics db.Characteristic `json:"characteristics"`
	DerivedStats    db.DerivedStat    `json:"derived_stats"`
	HP              db.HealthState    `json:"hp"`
	MP              db.MagicState     `json:"mp"`
	Sanity          db.SanityState    `json:"sanity"`
	Luck            db.LuckState      `json:"luck"`
	Backstory       BackstoryModel    `json:"backstory"`
	Finances        FinancesModel     `json:"finances"`
	Notes           db.Note           `json:"notes"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
