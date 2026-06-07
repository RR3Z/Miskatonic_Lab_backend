package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
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
