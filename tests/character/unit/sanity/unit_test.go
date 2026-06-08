package tests

import (
	"context"
	"errors"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestUpsertSanityAllowsCurrentSanityLessThanMaxSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(sanityInt16(60), sanityInt16(40))
	expectedSanity := testSanityState()
	dbtx.QueryRowData = sanityRowData(expectedSanity)

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityAllowsCurrentSanityEqualToMaxSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(sanityInt16(60), sanityInt16(60))
	expectedSanity := testSanityState()
	expectedSanity.CurrentSanity = 60
	dbtx.QueryRowData = sanityRowData(expectedSanity)

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityRejectsCurrentSanityGreaterThanMaxSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()

	_, err := characterService.UpsertSanity(context.Background(), testUpsertSanityInput(sanityInt16(5), sanityInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentSanityExceedsMax)
	require.Equal(t, 0, dbtx.QueryRowCalls)
}

func TestUpsertSanityUsesExistingStateForPartialValidation(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	existingSanity := testSanityState()
	existingSanity.MaxSanity = 5
	existingSanity.CurrentSanity = 4
	dbtx.QueryRowData = sanityRowData(existingSanity)

	_, err := characterService.UpsertSanity(context.Background(), testUpsertSanityInput(nil, sanityInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentSanityExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertSanityRejectsPartialMaxSanityBelowExistingCurrentSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	existingSanity := testSanityState()
	existingSanity.MaxSanity = 10
	existingSanity.CurrentSanity = 7
	dbtx.QueryRowData = sanityRowData(existingSanity)

	_, err := characterService.UpsertSanity(context.Background(), testUpsertSanityInput(sanityInt16(5), nil))

	require.ErrorIs(t, err, myErrors.ErrCurrentSanityExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertSanityAllowsPartialCurrentSanityWithinExistingMaxSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(nil, sanityInt16(6))
	existingSanity := testSanityState()
	existingSanity.MaxSanity = 10
	existingSanity.CurrentSanity = 4
	expectedSanity := testSanityState()
	expectedSanity.MaxSanity = 10
	expectedSanity.CurrentSanity = 6
	dbtx.QueryRowResults = []FakeSanityQueryRowResult{
		{Data: sanityRowData(existingSanity)},
		{Data: sanityRowData(expectedSanity)},
	}

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityAllowsPartialMaxSanityAboveExistingCurrentSanity(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(sanityInt16(10), nil)
	existingSanity := testSanityState()
	existingSanity.MaxSanity = 5
	existingSanity.CurrentSanity = 7
	expectedSanity := testSanityState()
	expectedSanity.MaxSanity = 10
	expectedSanity.CurrentSanity = 7
	dbtx.QueryRowResults = []FakeSanityQueryRowResult{
		{Data: sanityRowData(existingSanity)},
		{Data: sanityRowData(expectedSanity)},
	}

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityReturnsExistingStateReadError(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	expectedErr := errors.New("get existing sanity failed")
	dbtx.QueryRowResults = []FakeSanityQueryRowResult{{Err: expectedErr}}

	_, err := characterService.UpsertSanity(context.Background(), testUpsertSanityInput(nil, sanityInt16(6)))

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertSanityAllowsNilNumericInputWithoutReadingExistingState(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(nil, nil)
	expectedSanity := testSanityState()
	dbtx.QueryRowData = sanityRowData(expectedSanity)

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	require.Equal(t, []any{
		input.UserID,
		input.CharacterID,
		input.MaxSanity,
		input.CurrentSanity,
		input.TempInsanity,
		input.IndefInsanity,
	}, dbtx.LastQueryRowArgs)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityAllowsBoolOnlyInputWithoutReadingExistingState(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(nil, nil)
	input.TempInsanity = sanityBool(true)
	expectedSanity := testSanityState()
	expectedSanity.TempInsanity = true
	dbtx.QueryRowData = sanityRowData(expectedSanity)

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	require.Equal(t, []any{
		input.UserID,
		input.CharacterID,
		input.MaxSanity,
		input.CurrentSanity,
		input.TempInsanity,
		input.IndefInsanity,
	}, dbtx.LastQueryRowArgs)
	requireSameSanityState(t, expectedSanity, sanity)
}

func TestUpsertSanityAllowsPartialInputWhenExistingStateIsMissing(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForSanity()
	input := testUpsertSanityInput(nil, sanityInt16(6))
	expectedSanity := testSanityState()
	expectedSanity.CurrentSanity = 6
	dbtx.QueryRowResults = []FakeSanityQueryRowResult{
		{Err: pgx.ErrNoRows},
		{Data: sanityRowData(expectedSanity)},
	}

	sanity, err := characterService.UpsertSanity(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	require.Equal(t, []any{
		input.UserID,
		input.CharacterID,
		input.MaxSanity,
		input.CurrentSanity,
		input.TempInsanity,
		input.IndefInsanity,
	}, dbtx.LastQueryRowArgs)
	requireSameSanityState(t, expectedSanity, sanity)
}
