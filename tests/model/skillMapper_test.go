package tests

import (
	"testing"

	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestToSkillModelCopiesProtectionAndBaseRule(t *testing.T) {
	row := testSkillRow()

	result := skillsDTO.ToSkillModel(row)

	requireSameSkill(t, row, result)
}

func TestToCharacterSkillModelsMapsAllRows(t *testing.T) {
	first := testSkillRow()
	second := testSkillRow()
	second.ID = testUUID("66666666-6666-6666-6666-666666666666")
	second.Name = "Science: Astronomy"
	second.IsProtected = false
	second.BaseRule = nil

	rows := []db.GetCharacterSkillsRow{
		characterSkillsRowFromGetSkillsRow(first),
		characterSkillsRowFromGetSkillsRow(second),
	}

	result := skillsDTO.ToCharacterSkillModels(rows)

	require.Len(t, result, 2)
	require.True(t, result[0].IsProtected)
	require.Equal(t, "dodge", *result[0].BaseRule)
	require.False(t, result[1].IsProtected)
	require.Nil(t, result[1].BaseRule)
	require.Equal(t, "Science: Astronomy", result[1].Name)
}

func TestToCharacterSkillModelsReturnsEmptySliceForEmptyInput(t *testing.T) {
	result := skillsDTO.ToCharacterSkillModels(nil)

	require.Empty(t, result)
	require.NotNil(t, result)
}

func TestSingleCreatedAndUpdatedSkillMappersKeepNewFields(t *testing.T) {
	row := testSkillRow()

	single := skillsDTO.ToSingleCharacterSkillModel(singleCharacterSkillRowFromGetSkillsRow(row))
	created := skillsDTO.ToCreatedCharacterSkillModel(createdCharacterSkillRowFromGetSkillsRow(row))
	updated := skillsDTO.ToUpdatedCharacterSkillModel(updatedCharacterSkillRowFromGetSkillsRow(row))

	for _, result := range []skillsDTO.SkillModel{single, created, updated} {
		require.Equal(t, row.Name, result.Name)
		require.True(t, result.IsProtected)
		require.Equal(t, "dodge", *result.BaseRule)
	}
}

func characterSkillsRowFromGetSkillsRow(row db.GetSkillsRow) db.GetCharacterSkillsRow {
	return db.GetCharacterSkillsRow{
		ID:           row.ID,
		CharacterID:  row.CharacterID,
		Name:         row.Name,
		CategoryID:   row.CategoryID,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CategoryName: row.CategoryName,
	}
}

func singleCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.GetCharacterSkillRow {
	return db.GetCharacterSkillRow{
		ID:           row.ID,
		CharacterID:  row.CharacterID,
		Name:         row.Name,
		CategoryID:   row.CategoryID,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CategoryName: row.CategoryName,
	}
}

func createdCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.CreateCharacterSkillRow {
	return db.CreateCharacterSkillRow{
		ID:           row.ID,
		CharacterID:  row.CharacterID,
		Name:         row.Name,
		CategoryID:   row.CategoryID,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CategoryName: row.CategoryName,
	}
}

func updatedCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.UpdateCharacterSkillRow {
	return db.UpdateCharacterSkillRow{
		ID:           row.ID,
		CharacterID:  row.CharacterID,
		Name:         row.Name,
		CategoryID:   row.CategoryID,
		BaseValue:    row.BaseValue,
		Value:        row.Value,
		Checked:      row.Checked,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		IsProtected:  row.IsProtected,
		BaseRule:     row.BaseRule,
		CategoryName: row.CategoryName,
	}
}
