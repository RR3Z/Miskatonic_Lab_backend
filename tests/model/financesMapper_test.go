package tests

import (
	"testing"

	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/stretchr/testify/require"
)

func TestToFinancesModelCopiesFieldsWithoutCreditRating(t *testing.T) {
	finance := testFinance()

	result := financesDTO.ToFinancesModel(finance, nil)

	require.Equal(t, finance.ID, result.ID)
	require.Equal(t, finance.SpendingLimit, result.SpendingLimit)
	require.Equal(t, finance.Cash, result.Cash)
	require.Equal(t, finance.Assets, result.Assets)
	require.Nil(t, result.CreditRating)
	require.Equal(t, finance.CreatedAt.Time, result.CreatedAt.Time)
	require.Equal(t, finance.UpdatedAt.Time, result.UpdatedAt.Time)
}

func TestToFinancesModelAttachesCreditRatingSkill(t *testing.T) {
	finance := testFinance()
	creditRating := skillsDTO.ToSkillModel(testSkillRow())

	result := financesDTO.ToFinancesModel(finance, &creditRating)

	require.NotNil(t, result.CreditRating)
	require.Equal(t, creditRating.ID, result.CreditRating.ID)
	require.Equal(t, creditRating.Name, result.CreditRating.Name)
}
