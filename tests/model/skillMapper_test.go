package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestToSkillModelCopiesSkillWithoutSpecialty(t *testing.T) {
	row := testSkillRow()

	result := model.ToSkillModel(row)

	requireSameSkill(t, row, result)
	require.Nil(t, result.Specialty)
}

func TestToSkillModelMapsSpecialtyWhenSpecialtyIDIsValid(t *testing.T) {
	row := testSpecializedSkillRow()

	result := model.ToSkillModel(row)

	requireSameSkill(t, row, result)
	require.NotNil(t, result.Specialty)
	require.Equal(t, row.SpecialtyPkID, result.Specialty.ID)
	require.Equal(t, *row.SpecialtyName, result.Specialty.Name)
	require.Equal(t, *row.SpecialtyDescription, result.Specialty.Description)
	require.Equal(t, *row.SpecialtyBaseValue, result.Specialty.BaseValue)
	require.Equal(t, row.SpecialtyCreatedAt.Time, result.Specialty.CreatedAt.Time)
	require.Equal(t, row.SpecialtyUpdatedAt.Time, result.Specialty.UpdatedAt.Time)
}

func TestToCharacterSkillModelsMapsAllRows(t *testing.T) {
	rows := []db.GetCharacterSkillsRow{
		testCharacterSkillsRowFromGetSkillsRow(testSkillRow()),
		testCharacterSkillsRowFromGetSkillsRow(testSpecializedSkillRow()),
	}

	result := model.ToCharacterSkillModels(rows)

	require.Len(t, result, 2)
	require.Equal(t, rows[0].ID, result[0].ID)
	require.Equal(t, rows[0].Name, result[0].Name)
	require.Nil(t, result[0].Specialty)
	require.Equal(t, rows[1].ID, result[1].ID)
	require.NotNil(t, result[1].Specialty)
	require.Equal(t, rows[1].SpecialtyPkID, result[1].Specialty.ID)
}

func TestToCharacterSkillModelsReturnsEmptySliceForEmptyInput(t *testing.T) {
	result := model.ToCharacterSkillModels(nil)

	require.Empty(t, result)
	require.NotNil(t, result)
}

func TestToSingleCharacterSkillModelMapsSpecialty(t *testing.T) {
	row := testSingleCharacterSkillRowFromGetSkillsRow(testSpecializedSkillRow())

	result := model.ToSingleCharacterSkillModel(row)

	require.Equal(t, row.ID, result.ID)
	require.Equal(t, row.Name, result.Name)
	require.Equal(t, row.BaseValue, result.BaseValue)
	require.Equal(t, row.Value, result.Value)
	require.Equal(t, row.Checked, result.Checked)
	require.Equal(t, row.CategoryName, result.Category)
	require.True(t, result.Specialized)
	require.NotNil(t, result.Specialty)
	require.Equal(t, row.SpecialtyPkID, result.Specialty.ID)
}

func TestToCreatedCharacterSkillModelMapsSpecialty(t *testing.T) {
	row := testCreatedCharacterSkillRowFromGetSkillsRow(testSpecializedSkillRow())

	result := model.ToCreatedCharacterSkillModel(row)

	require.Equal(t, row.ID, result.ID)
	require.Equal(t, row.Name, result.Name)
	require.Equal(t, row.BaseValue, result.BaseValue)
	require.Equal(t, row.Value, result.Value)
	require.Equal(t, row.Checked, result.Checked)
	require.Equal(t, row.CategoryName, result.Category)
	require.True(t, result.Specialized)
	require.NotNil(t, result.Specialty)
	require.Equal(t, row.SpecialtyPkID, result.Specialty.ID)
}

func TestToUpdatedCharacterSkillModelMapsSpecialty(t *testing.T) {
	row := testUpdatedCharacterSkillRowFromGetSkillsRow(testSpecializedSkillRow())

	result := model.ToUpdatedCharacterSkillModel(row)

	require.Equal(t, row.ID, result.ID)
	require.Equal(t, row.Name, result.Name)
	require.Equal(t, row.BaseValue, result.BaseValue)
	require.Equal(t, row.Value, result.Value)
	require.Equal(t, row.Checked, result.Checked)
	require.Equal(t, row.CategoryName, result.Category)
	require.True(t, result.Specialized)
	require.NotNil(t, result.Specialty)
	require.Equal(t, row.SpecialtyPkID, result.Specialty.ID)
}

func testCharacterSkillsRowFromGetSkillsRow(row db.GetSkillsRow) db.GetCharacterSkillsRow {
	return db.GetCharacterSkillsRow{
		ID:                   row.ID,
		CharacterID:          row.CharacterID,
		Name:                 row.Name,
		CategoryID:           row.CategoryID,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		Specialized:          row.Specialized,
		SpecialtyID:          row.SpecialtyID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CategoryName:         row.CategoryName,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
	}
}

func testSingleCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.GetCharacterSkillRow {
	return db.GetCharacterSkillRow{
		ID:                   row.ID,
		CharacterID:          row.CharacterID,
		Name:                 row.Name,
		CategoryID:           row.CategoryID,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		Specialized:          row.Specialized,
		SpecialtyID:          row.SpecialtyID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CategoryName:         row.CategoryName,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
	}
}

func testCreatedCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.CreateCharacterSkillRow {
	return db.CreateCharacterSkillRow{
		ID:                   row.ID,
		CharacterID:          row.CharacterID,
		Name:                 row.Name,
		CategoryID:           row.CategoryID,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		Specialized:          row.Specialized,
		SpecialtyID:          row.SpecialtyID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CategoryName:         row.CategoryName,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
	}
}

func testUpdatedCharacterSkillRowFromGetSkillsRow(row db.GetSkillsRow) db.UpdateCharacterSkillRow {
	return db.UpdateCharacterSkillRow{
		ID:                   row.ID,
		CharacterID:          row.CharacterID,
		Name:                 row.Name,
		CategoryID:           row.CategoryID,
		BaseValue:            row.BaseValue,
		Value:                row.Value,
		Checked:              row.Checked,
		Specialized:          row.Specialized,
		SpecialtyID:          row.SpecialtyID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CategoryName:         row.CategoryName,
		SpecialtyPkID:        row.SpecialtyPkID,
		SpecialtyName:        row.SpecialtyName,
		SpecialtyDescription: row.SpecialtyDescription,
		SpecialtyBaseValue:   row.SpecialtyBaseValue,
		SpecialtyCreatedAt:   row.SpecialtyCreatedAt,
		SpecialtyUpdatedAt:   row.SpecialtyUpdatedAt,
	}
}
