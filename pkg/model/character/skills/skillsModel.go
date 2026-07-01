package skills

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SkillSpecialtyModel struct {
	ID pgtype.UUID `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	BaseValue   int16  `json:"base_value"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

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

type skillModelFields struct {
	ID pgtype.UUID

	Name         string
	BaseValue    int16
	Value        int16
	Checked      bool
	CategoryName string
	Specialized  bool

	SpecialtyPkID        pgtype.UUID
	SpecialtyName        *string
	SpecialtyDescription *string
	SpecialtyBaseValue   *int16
	SpecialtyCreatedAt   pgtype.Timestamptz
	SpecialtyUpdatedAt   pgtype.Timestamptz

	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func skillModelFromFields(row skillModelFields) SkillModel {
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

func ToSkillModel(row db.GetSkillsRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:                   row.ID,
		Name:                 row.Name,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		CategoryName:         row.CategoryName,
		Specialized:          row.Specialized,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	})
}

func ToCharacterSkillModels(rows []db.GetCharacterSkillsRow) []SkillModel {
	models := make([]SkillModel, len(rows))
	for i, row := range rows {
		models[i] = ToCharacterSkillModel(row)
	}
	return models
}

func ToCharacterSkillModel(row db.GetCharacterSkillsRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:                   row.ID,
		Name:                 row.Name,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		CategoryName:         row.CategoryName,
		Specialized:          row.Specialized,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	})
}

func ToSingleCharacterSkillModel(row db.GetCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:                   row.ID,
		Name:                 row.Name,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		CategoryName:         row.CategoryName,
		Specialized:          row.Specialized,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	})
}

func ToCreatedCharacterSkillModel(row db.CreateCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:                   row.ID,
		Name:                 row.Name,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		CategoryName:         row.CategoryName,
		Specialized:          row.Specialized,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	})
}

func ToUpdatedCharacterSkillModel(row db.UpdateCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:                   row.ID,
		Name:                 row.Name,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		CategoryName:         row.CategoryName,
		Specialized:          row.Specialized,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	})
}
