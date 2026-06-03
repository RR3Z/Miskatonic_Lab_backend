package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SkillModel struct {
	ID pgtype.UUID `json:"id"`

	Name        string               `json:"name"`
	BaseValue   int16                `json:"base_value"`
	Value       int16                `json:"value"`
	Checked     bool                 `json:"checked"`
	Category    string               `json:"category"`
	Specialized bool                 `json:"specialized"`
	Specialty   *SkillSpecialtyModel `json:"specialty"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type SkillSpecialtyModel struct {
	ID pgtype.UUID `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	BaseValue   int16  `json:"base_value"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
