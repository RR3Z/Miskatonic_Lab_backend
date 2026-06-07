package tests

import (
	"context"
	"errors"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestUpsertHealthAllowsCurrentHpLessThanMaxHp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	input := testUpsertHealthInput(healthInt16(10), healthInt16(7))
	expectedHealth := testHealthState()
	dbtx.QueryRowData = healthRowData(expectedHealth)

	health, err := characterService.UpsertHealth(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameHealthState(t, expectedHealth, health)
}

func TestUpsertHealthAllowsCurrentHpEqualToMaxHp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	input := testUpsertHealthInput(healthInt16(10), healthInt16(10))
	expectedHealth := testHealthState()
	expectedHealth.CurrentHp = 10
	dbtx.QueryRowData = healthRowData(expectedHealth)

	health, err := characterService.UpsertHealth(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, dbtx.QueryRowCalls)
	requireSameHealthState(t, expectedHealth, health)
}

func TestUpsertHealthRejectsCurrentHpGreaterThanMaxHp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()

	_, err := characterService.UpsertHealth(context.Background(), testUpsertHealthInput(healthInt16(5), healthInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentHealthExceedsMax)
	require.Equal(t, 0, dbtx.QueryRowCalls)
}

func TestUpsertHealthUsesExistingStateForPartialValidation(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	existingHealth := testHealthState()
	existingHealth.MaxHp = 5
	existingHealth.CurrentHp = 4
	dbtx.QueryRowData = healthRowData(existingHealth)

	_, err := characterService.UpsertHealth(context.Background(), testUpsertHealthInput(nil, healthInt16(6)))

	require.ErrorIs(t, err, myErrors.ErrCurrentHealthExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertHealthRejectsPartialMaxHpBelowExistingCurrentHp(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	existingHealth := testHealthState()
	existingHealth.MaxHp = 10
	existingHealth.CurrentHp = 7
	dbtx.QueryRowData = healthRowData(existingHealth)

	_, err := characterService.UpsertHealth(context.Background(), testUpsertHealthInput(healthInt16(5), nil))

	require.ErrorIs(t, err, myErrors.ErrCurrentHealthExceedsMax)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertHealthReturnsExistingStateReadError(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	expectedErr := errors.New("get existing health failed")
	dbtx.QueryRowResults = []FakeHealthQueryRowResult{{Err: expectedErr}}

	_, err := characterService.UpsertHealth(context.Background(), testUpsertHealthInput(nil, healthInt16(6)))

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 1, dbtx.QueryRowCalls)
}

func TestUpsertHealthAllowsPartialInputWhenExistingStateIsMissing(t *testing.T) {
	dbtx, characterService := newTestCharacterServiceForHealth()
	input := testUpsertHealthInput(nil, healthInt16(6))
	expectedHealth := testHealthState()
	expectedHealth.CurrentHp = 6
	dbtx.QueryRowResults = []FakeHealthQueryRowResult{
		{Err: pgx.ErrNoRows},
		{Data: healthRowData(expectedHealth)},
	}

	health, err := characterService.UpsertHealth(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 2, dbtx.QueryRowCalls)
	require.Equal(t, []any{
		input.UserID,
		input.CharacterID,
		input.MaxHp,
		input.CurrentHp,
		input.MajorWound,
		input.Unconscious,
		input.Dying,
		input.Dead,
	}, dbtx.LastQueryRowArgs)
	requireSameHealthState(t, expectedHealth, health)
}
