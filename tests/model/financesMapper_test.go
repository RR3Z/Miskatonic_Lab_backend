package tests

import (
	"testing"

	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/stretchr/testify/require"
)

func TestToFinancesModelCopiesFields(t *testing.T) {
	finance := testFinance()

	result := financesDTO.ToFinancesModel(finance)

	require.Equal(t, finance.ID, result.ID)
	require.Equal(t, finance.SpendingLimit, result.SpendingLimit)
	require.Equal(t, finance.Cash, result.Cash)
	require.Equal(t, finance.Assets, result.Assets)
	require.Equal(t, finance.CreatedAt.Time, result.CreatedAt.Time)
	require.Equal(t, finance.UpdatedAt.Time, result.UpdatedAt.Time)
}
