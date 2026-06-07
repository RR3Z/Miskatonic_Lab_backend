package tests

import (
	"context"
	"errors"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestUpsertLuckAllowsCurrentLuckLessThanStartingLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(luckInt16(60), luckInt16(40))
	expectedLuck := testLuckState()
	dbtx.QueryRowData = luckRowData(expectedLuck)

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameLuckState(t, expectedLuck, luck)
}

func TestUpsertLuckAllowsCurrentLuckEqualToStartingLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(luckInt16(60), luckInt16(60))
	expectedLuck := testLuckState()
	expectedLuck.CurrentLuck = 60
	dbtx.QueryRowData = luckRowData(expectedLuck)

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameLuckState(t, expectedLuck, luck)
}

func TestUpsertLuckRejectsCurrentLuckGreaterThanStartingLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()

	_, err := characterService.UpsertLuck(context.Background(), testUpsertLuckInput(luckInt16(5), luckInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentLuckExceedsStarting)
	require.Equal(t, 0, dbtx.QueryRowCalls)
}

func TestUpsertLuckUsesExistingStateForPartialValidation(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	existingLuck := testLuckState()
	existingLuck.StartingLuck = 5
	existingLuck.CurrentLuck = 4
	dbtx.QueryRowData = luckRowData(existingLuck)

	_, err := characterService.UpsertLuck(context.Background(), testUpsertLuckInput(nil, luckInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentLuckExceedsStarting)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertLuckRejectsPartialStartingLuckBelowExistingCurrentLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	existingLuck := testLuckState()
	existingLuck.StartingLuck = 10
	existingLuck.CurrentLuck = 7
	dbtx.QueryRowData = luckRowData(existingLuck)

	_, err := characterService.UpsertLuck(context.Background(), testUpsertLuckInput(luckInt16(5), nil))

	require.ErrorIs(t, err, myErrors.ErrCurrentLuckExceedsStarting)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertLuckReturnsExistingStateReadError(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	expectedErr := errors.New("get existing luck failed")
	dbtx.QueryRowResults = []FakeLuckQueryRowResult{{Err: expectedErr}}

	_, err := characterService.UpsertLuck(context.Background(), testUpsertLuckInput(nil, luckInt16(6)))

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertLuckAllowsPartialInputWhenExistingStateIsMissing(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(nil, luckInt16(6))
	expectedLuck := testLuckState()
	expectedLuck.CurrentLuck = 6
	dbtx.QueryRowResults = []FakeLuckQueryRowResult{
		{Err: pgx.ErrNoRows},
		{Data: luckRowData(expectedLuck)},
	}

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	require.Equal(t, []any{input.UserID, input.CharacterID, input.StartingLuck, input.CurrentLuck}, dbtx.LastQueryRowArgs)
	requireSameLuckState(t, expectedLuck, luck)
}
