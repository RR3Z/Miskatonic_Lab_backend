package model

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
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

func ToSkillModel(row db.GetSkillsRow) SkillModel {
	var specialty *SkillSpecialtyModel
	if row.SpecialtyPkID.Valid {
		specialty = &SkillSpecialtyModel{
			ID:          row.SpecialtyPkID,
			Name:        *row.SpecialtyName,
			Description: *row.SpecialtyDescription,
			BaseValue:   *row.SpecialtyBaseValue,
			CreatedAt:   row.SpecialtyCreatedAt,
			UpdatedAt:   row.SpecialtyUpdatedAt,
		}
	}

	return SkillModel{
		ID:          row.ID,
		Name:        row.Name,
		BaseValue:   row.BaseValue,
		Value:       row.Value,
		Checked:     row.Checked,
		Category:    row.CategoryName,
		Specialized: row.Specialized,
		Specialty:   specialty,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
