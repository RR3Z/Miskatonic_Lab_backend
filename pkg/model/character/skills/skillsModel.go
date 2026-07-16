package skillsDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SkillModel struct {
	ID pgtype.UUID `json:"id"`

	Name        string  `json:"name"`
	BaseValue   int16   `json:"base_value"`
	Value       int16   `json:"value"`
	Checked     bool    `json:"checked"`
	Category    string  `json:"category"`
	IsProtected bool    `json:"is_protected"`
	BaseRule    *string `json:"base_rule"`

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
	IsProtected  bool
	BaseRule     *string

	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func skillModelFromFields(row skillModelFields) SkillModel {
	return SkillModel{
		ID:          row.ID,
		Name:        row.Name,
		BaseValue:   row.BaseValue,
		Value:       row.Value,
		Checked:     row.Checked,
		Category:    row.CategoryName,
		IsProtected: row.IsProtected,
		BaseRule:    row.BaseRule,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func ToSkillModel(row db.GetSkillsRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:           row.ID,
		Name:         row.Name,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CategoryName: row.CategoryName,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
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
		ID:           row.ID,
		Name:         row.Name,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CategoryName: row.CategoryName,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	})
}

func ToSingleCharacterSkillModel(row db.GetCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:           row.ID,
		Name:         row.Name,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CategoryName: row.CategoryName,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	})
}

func ToCreatedCharacterSkillModel(row db.CreateCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:           row.ID,
		Name:         row.Name,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CategoryName: row.CategoryName,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	})
}

func ToUpdatedCharacterSkillModel(row db.UpdateCharacterSkillRow) SkillModel {
	return skillModelFromFields(skillModelFields{
		ID:           row.ID,
		Name:         row.Name,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CategoryName: row.CategoryName,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	})
}
