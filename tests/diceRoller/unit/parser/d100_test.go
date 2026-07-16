package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/dice"
	"github.com/stretchr/testify/require"
)

func TestBuildD100CandidatesSharesUnitsAcrossTensDice(t *testing.T) {
	candidates, err := dice.BuildD100Candidates(4, []int{2, 4})

	require.NoError(t, err)
	require.Equal(t, []int{24, 44}, candidates)
}

func TestBuildD100CandidatesTreatsDoubleZeroAsOneHundred(t *testing.T) {
	candidates, err := dice.BuildD100Candidates(0, []int{0, 1})

	require.NoError(t, err)
	require.Equal(t, []int{100, 10}, candidates)
}

func TestBuildD100CandidatesRejectsInvalidDigits(t *testing.T) {
	_, err := dice.BuildD100Candidates(10, []int{0})
	require.Error(t, err)

	_, err = dice.BuildD100Candidates(0, []int{10})
	require.Error(t, err)

	_, err = dice.BuildD100Candidates(0, nil)
	require.Error(t, err)
}
