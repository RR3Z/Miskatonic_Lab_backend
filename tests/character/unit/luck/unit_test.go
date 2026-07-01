package tests

import (
	"context"
	"errors"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
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

func TestUpsertLuckAllowsPartialCurrentLuckWithinExistingStartingLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(nil, luckInt16(6))
	existingLuck := testLuckState()
	existingLuck.StartingLuck = 10
	existingLuck.CurrentLuck = 4
	expectedLuck := testLuckState()
	expectedLuck.StartingLuck = 10
	expectedLuck.CurrentLuck = 6
	dbtx.QueryRowResults = []FakeLuckQueryRowResult{
		{Data: luckRowData(existingLuck)},
		{Data: luckRowData(expectedLuck)},
	}

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameLuckState(t, expectedLuck, luck)
}

func TestUpsertLuckAllowsPartialStartingLuckAboveExistingCurrentLuck(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(luckInt16(10), nil)
	existingLuck := testLuckState()
	existingLuck.StartingLuck = 5
	existingLuck.CurrentLuck = 7
	expectedLuck := testLuckState()
	expectedLuck.StartingLuck = 10
	expectedLuck.CurrentLuck = 7
	dbtx.QueryRowResults = []FakeLuckQueryRowResult{
		{Data: luckRowData(existingLuck)},
		{Data: luckRowData(expectedLuck)},
	}

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameLuckState(t, expectedLuck, luck)
}

func TestUpsertLuckReturnsExistingStateReadError(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	expectedErr := errors.New("get existing luck failed")
	dbtx.QueryRowResults = []FakeLuckQueryRowResult{{Err: expectedErr}}

	_, err := characterService.UpsertLuck(context.Background(), testUpsertLuckInput(nil, luckInt16(6)))

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertLuckAllowsNilNumericInputWithoutReadingExistingState(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForLuck()
	input := testUpsertLuckInput(nil, nil)
	expectedLuck := testLuckState()
	dbtx.QueryRowData = luckRowData(expectedLuck)

	luck, err := characterService.UpsertLuck(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	require.Equal(t, []any{input.UserID, input.CharacterID, input.StartingLuck, input.CurrentLuck}, dbtx.LastQueryRowArgs)
	requireSameLuckState(t, expectedLuck, luck)
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

func TestUpsertLuckRejectsNegativeStartingLuck(t *testing.T) {
	_, service := newTestCharacterServiceForLuck()
	startingLuck := int16(-5)
	_, err := service.UpsertLuck(context.Background(), testUpsertLuckInput(&startingLuck, nil))
	require.ErrorIs(t, err, characterErrors.ErrStateNegative)
}

func TestUpsertLuckRejectsNegativeCurrentLuck(t *testing.T) {
	_, service := newTestCharacterServiceForLuck()
	currentLuck := int16(-3)
	_, err := service.UpsertLuck(context.Background(), testUpsertLuckInput(nil, &currentLuck))
	require.ErrorIs(t, err, characterErrors.ErrStateNegative)
}
