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

func TestUpsertMagicAllowsCurrentMpLessThanMaxMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	input := testUpsertMagicInput(magicInt16(10), magicInt16(7))
	expectedMagic := testMagicState()
	dbtx.QueryRowData = magicRowData(expectedMagic)

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameMagicState(t, expectedMagic, magic)
}

func TestUpsertMagicAllowsCurrentMpEqualToMaxMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	input := testUpsertMagicInput(magicInt16(10), magicInt16(10))
	expectedMagic := testMagicState()
	expectedMagic.CurrentMp = 10
	dbtx.QueryRowData = magicRowData(expectedMagic)

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameMagicState(t, expectedMagic, magic)
}

// INPUT -> MaxMP = 5, CurrentMP = 6 -> Error
func TestUpsertMagicRejectsCurrentMpGreaterThanMaxMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()

	_, err := characterService.UpsertMagic(context.Background(), testUpsertMagicInput(magicInt16(5), magicInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentMagicExceedsMax)
	require.Equal(t, 0, dbtx.QueryRowCalls)
}

// EXISTING_DATA -> MaxMP = 5, CurrentMP = 4
// INPUT -> MaxMP = nil, CurrentMP = 6 -> Error
func TestUpsertMagicUsesExistingStateForPartialValidation(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()

	existingMagic := testMagicState()
	existingMagic.MaxMp = 5
	existingMagic.CurrentMp = 4

	dbtx.QueryRowData = magicRowData(existingMagic)

	_, err := characterService.UpsertMagic(context.Background(), testUpsertMagicInput(nil, magicInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentMagicExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertMagicRejectsPartialMaxMpBelowExistingCurrentMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	existingMagic := testMagicState()
	existingMagic.MaxMp = 10
	existingMagic.CurrentMp = 7
	dbtx.QueryRowData = magicRowData(existingMagic)

	_, err := characterService.UpsertMagic(context.Background(), testUpsertMagicInput(magicInt16(5), nil))

	require.ErrorIs(t, err, myErrors.ErrCurrentMagicExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertMagicAllowsPartialCurrentMpWithinExistingMaxMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	input := testUpsertMagicInput(nil, magicInt16(6))
	existingMagic := testMagicState()
	existingMagic.MaxMp = 10
	existingMagic.CurrentMp = 4
	expectedMagic := testMagicState()
	expectedMagic.MaxMp = 10
	expectedMagic.CurrentMp = 6
	dbtx.QueryRowResults = []FakeMagicQueryRowResult{
		{Data: magicRowData(existingMagic)},
		{Data: magicRowData(expectedMagic)},
	}

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameMagicState(t, expectedMagic, magic)
}

func TestUpsertMagicAllowsPartialMaxMpAboveExistingCurrentMp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	input := testUpsertMagicInput(magicInt16(10), nil)
	existingMagic := testMagicState()
	existingMagic.MaxMp = 5
	existingMagic.CurrentMp = 7
	expectedMagic := testMagicState()
	expectedMagic.MaxMp = 10
	expectedMagic.CurrentMp = 7
	dbtx.QueryRowResults = []FakeMagicQueryRowResult{
		{Data: magicRowData(existingMagic)},
		{Data: magicRowData(expectedMagic)},
	}

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameMagicState(t, expectedMagic, magic)
}

func TestUpsertMagicReturnsExistingStateReadError(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	expectedErr := errors.New("get existing magic failed")
	dbtx.QueryRowResults = []FakeMagicQueryRowResult{{Err: expectedErr}}

	_, err := characterService.UpsertMagic(context.Background(), testUpsertMagicInput(nil, magicInt16(6)))

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertMagicAllowsNilNumericInputWithoutReadingExistingState(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()
	input := testUpsertMagicInput(nil, nil)
	expectedMagic := testMagicState()
	dbtx.QueryRowData = magicRowData(expectedMagic)

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	require.Equal(t, []any{input.UserID, input.CharacterID, input.MaxMp, input.CurrentMp}, dbtx.LastQueryRowArgs)
	requireSameMagicState(t, expectedMagic, magic)
}

// EXISTING_DATA -> MaxMP = nil, CurrentMP = 6
// INPUT -> MaxMP = nil, CurrentMP = 6 -> MagicState
func TestUpsertMagicAllowsPartialInputWhenExistingStateIsMissing(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForMagic()

	input := testUpsertMagicInput(nil, magicInt16(6))

	expectedMagic := testMagicState()
	expectedMagic.CurrentMp = 6
	dbtx.QueryRowResults = []FakeMagicQueryRowResult{
		{Err: pgx.ErrNoRows},
		{Data: magicRowData(expectedMagic)},
	}

	magic, err := characterService.UpsertMagic(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	require.Equal(t, []any{input.UserID, input.CharacterID, input.MaxMp, input.CurrentMp}, dbtx.LastQueryRowArgs)
	requireSameMagicState(t, expectedMagic, magic)
}

func TestUpsertMagicRejectsNegativeMaxMp(t *testing.T) {
	_, service := newTestCharacterServiceForMagic()
	maxMp := int16(-5)
	_, err := service.UpsertMagic(context.Background(), testUpsertMagicInput(&maxMp, nil))
	require.ErrorIs(t, err, characterErrors.ErrStateNegative)
}

func TestUpsertMagicRejectsNegativeCurrentMp(t *testing.T) {
	_, service := newTestCharacterServiceForMagic()
	currentMp := int16(-3)
	_, err := service.UpsertMagic(context.Background(), testUpsertMagicInput(nil, &currentMp))
	require.ErrorIs(t, err, characterErrors.ErrStateNegative)
}
