package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestFinancesTableUpsertCreatesGetsAndPartiallyUpdatesFinances(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
		Cash:          financeString("$120"),
		Assets:        financeString("A battered motorcar"),
	})
	require.NoError(t, err)

	require.True(t, createdFinances.ID.Valid)
	require.Equal(t, character.ID, createdFinances.CharacterID)
	requireFinanceString(t, createdFinances.SpendingLimit, "$50")
	requireFinanceString(t, createdFinances.Cash, "$120")
	requireFinanceString(t, createdFinances.Assets, "A battered motorcar")
	require.True(t, createdFinances.CreatedAt.Valid)
	require.True(t, createdFinances.UpdatedAt.Valid)

	fetchedFinances, err := subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdFinances.ID, fetchedFinances.ID)

	time.Sleep(5 * time.Millisecond)

	updatedFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$75"),
		Cash:          financeString("$200"),
	})
	require.NoError(t, err)

	require.Equal(t, createdFinances.ID, updatedFinances.ID)
	requireFinanceString(t, updatedFinances.SpendingLimit, "$75")
	requireFinanceString(t, updatedFinances.Cash, "$200")
	requireFinanceString(t, updatedFinances.Assets, "A battered motorcar")
	require.True(t, updatedFinances.UpdatedAt.Time.After(createdFinances.UpdatedAt.Time) || updatedFinances.UpdatedAt.Time.Equal(createdFinances.UpdatedAt.Time))
}

func TestFinancesTableUpsertAllowsAllNilValuesOnInsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	finances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.True(t, finances.ID.Valid)
	require.Equal(t, character.ID, finances.CharacterID)
	require.Nil(t, finances.SpendingLimit)
	require.Nil(t, finances.Cash)
	require.Nil(t, finances.Assets)
}

func TestFinancesTableNilUpdateDoesNotOverwriteExistingValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
		Cash:          financeString("$120"),
		Assets:        financeString("A battered motorcar"),
	})
	require.NoError(t, err)

	updatedFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, createdFinances.ID, updatedFinances.ID)
	requireFinanceString(t, updatedFinances.SpendingLimit, "$50")
	requireFinanceString(t, updatedFinances.Cash, "$120")
	requireFinanceString(t, updatedFinances.Assets, "A battered motorcar")
}

func TestFinancesTablePartialUpdateAfterNilInsertOnlySetsProvidedValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	updatedFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$75"),
	})
	require.NoError(t, err)

	require.Equal(t, createdFinances.ID, updatedFinances.ID)
	requireFinanceString(t, updatedFinances.SpendingLimit, "$75")
	require.Nil(t, updatedFinances.Cash)
	require.Nil(t, updatedFinances.Assets)
}

func TestFinancesTableAllowsEmptyStrings(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	finances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString(""),
		Cash:          financeString(""),
		Assets:        financeString(""),
	})
	require.NoError(t, err)

	requireFinanceString(t, finances.SpendingLimit, "")
	requireFinanceString(t, finances.Cash, "")
	requireFinanceString(t, finances.Assets, "")
}

func TestFinancesTableAllowsBoundaryLengthMoneyFieldsAndLongAssets(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	boundaryValue := strings.Repeat("a", 120)
	longAssets := strings.Repeat("asset ", 1000)

	finances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString(boundaryValue),
		Cash:          financeString(boundaryValue),
		Assets:        financeString(longAssets),
	})
	require.NoError(t, err)

	requireFinanceString(t, finances.SpendingLimit, boundaryValue)
	requireFinanceString(t, finances.Cash, boundaryValue)
	requireFinanceString(t, finances.Assets, longAssets)
}

func TestFinancesTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        otherUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        owner.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
	})
	require.NoError(t, err)

	_, err = subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteFinances(context.Background(), db.DeleteFinancesParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedFinances, err := subject.queries.DeleteFinances(context.Background(), db.DeleteFinancesParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdFinances.ID, deletedFinances.ID)
}

func TestFinancesTableUnauthorizedUpsertDoesNotMutateExistingFinances(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        owner.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
		Cash:          financeString("$120"),
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        otherUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$1"),
		Cash:          financeString("$1"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedFinances, err := subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdFinances.ID, fetchedFinances.ID)
	requireFinanceString(t, fetchedFinances.SpendingLimit, "$50")
	requireFinanceString(t, fetchedFinances.Cash, "$120")
}

func TestFinancesTableReturnsNoRowsBeforeUpsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteFinances(context.Background(), db.DeleteFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestFinancesTableKeepsFinancesScopedToRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   firstCharacter.ID,
		SpendingLimit: financeString("$50"),
	})
	require.NoError(t, err)

	secondFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   secondCharacter.ID,
		SpendingLimit: financeString("$500"),
	})
	require.NoError(t, err)

	fetchedFirstFinances, err := subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstFinances.ID, fetchedFirstFinances.ID)
	requireFinanceString(t, fetchedFirstFinances.SpendingLimit, "$50")

	fetchedSecondFinances, err := subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, secondFinances.ID, fetchedSecondFinances.ID)
	requireFinanceString(t, fetchedSecondFinances.SpendingLimit, "$500")
}

func TestFinancesTableReturnsNoRowsForMissingCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	missingCharacterID := characterTestUUID("89898989-8989-8989-8989-898989898989")

	_, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   missingCharacterID,
		SpendingLimit: financeString("$50"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteFinances(context.Background(), db.DeleteFinancesParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestFinancesTableTruncatesTooLongMoneyFields(t *testing.T) {
	tests := []struct {
		name     string
		params   func(userID string, characterID pgtype.UUID, longValue string) db.UpsertFinancesParams
		assertFn func(t *testing.T, finances db.Finance, expectedValue string)
	}{
		{
			name: "spending limit",
			params: func(userID string, characterID pgtype.UUID, longValue string) db.UpsertFinancesParams {
				return db.UpsertFinancesParams{
					UserID:        userID,
					CharacterID:   characterID,
					SpendingLimit: financeString(longValue),
				}
			},
			assertFn: func(t *testing.T, finances db.Finance, expectedValue string) {
				requireFinanceString(t, finances.SpendingLimit, expectedValue)
			},
		},
		{
			name: "cash",
			params: func(userID string, characterID pgtype.UUID, longValue string) db.UpsertFinancesParams {
				return db.UpsertFinancesParams{
					UserID:      userID,
					CharacterID: characterID,
					Cash:        financeString(longValue),
				}
			},
			assertFn: func(t *testing.T, finances db.Finance, expectedValue string) {
				requireFinanceString(t, finances.Cash, expectedValue)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)
			longValue := strings.Repeat("a", 121)

			finances, err := subject.queries.UpsertFinances(context.Background(), tc.params(testUser.ID, character.ID, longValue))
			require.NoError(t, err)
			tc.assertFn(t, finances, strings.Repeat("a", 120))
		})
	}
}

func TestFinancesTableAllowsDeletingUnrelatedSkill(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	skill := createFinanceTestSkill(t, subject, testUser.ID, character.ID, "Средства")

	_, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.NoError(t, err)
}

func TestFinancesTableDeleteReturnsDeletedValuesAndAllowsRecreate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
	})
	require.NoError(t, err)

	deletedFinances, err := subject.queries.DeleteFinances(context.Background(), db.DeleteFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdFinances.ID, deletedFinances.ID)
	requireFinanceString(t, deletedFinances.SpendingLimit, "$50")

	recreatedFinances, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$75"),
	})
	require.NoError(t, err)
	require.NotEqual(t, deletedFinances.ID, recreatedFinances.ID)
	requireFinanceString(t, recreatedFinances.SpendingLimit, "$75")
}

func TestFinancesTableDeletingCharacterCascadesFinances(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.UpsertFinances(context.Background(), db.UpsertFinancesParams{
		UserID:        testUser.ID,
		CharacterID:   character.ID,
		SpendingLimit: financeString("$50"),
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetFinances(context.Background(), db.GetFinancesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
