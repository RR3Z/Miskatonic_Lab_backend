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

func TestBackstoriesTableUpsertCreatesGetsAndPartiallyUpdatesBackstory(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("A careful scholar with a dangerous curiosity."),
	})
	require.NoError(t, err)

	require.True(t, createdBackstory.ID.Valid)
	require.Equal(t, character.ID, createdBackstory.CharacterID)
	requireBackstoryString(t, createdBackstory.PersonalDescription, "A careful scholar with a dangerous curiosity.")
	require.True(t, createdBackstory.CreatedAt.Valid)
	require.True(t, createdBackstory.UpdatedAt.Valid)

	fetchedBackstory, err := subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdBackstory.ID, fetchedBackstory.ID)

	time.Sleep(5 * time.Millisecond)

	updatedBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("A retired professor who knows too much."),
	})
	require.NoError(t, err)

	require.Equal(t, createdBackstory.ID, updatedBackstory.ID)
	requireBackstoryString(t, updatedBackstory.PersonalDescription, "A retired professor who knows too much.")
	require.True(t, updatedBackstory.UpdatedAt.Time.After(createdBackstory.UpdatedAt.Time) || updatedBackstory.UpdatedAt.Time.Equal(createdBackstory.UpdatedAt.Time))
}

func TestBackstoriesTableUpsertAllowsNilEmptyAndLongPersonalDescription(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	nilCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	emptyCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	longCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	longDescription := strings.Repeat("description ", 1000)

	nilBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: nilCharacter.ID,
	})
	require.NoError(t, err)
	require.Nil(t, nilBackstory.PersonalDescription)

	emptyBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         emptyCharacter.ID,
		PersonalDescription: backstoryString(""),
	})
	require.NoError(t, err)
	requireBackstoryString(t, emptyBackstory.PersonalDescription, "")

	longBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         longCharacter.ID,
		PersonalDescription: backstoryString(longDescription),
	})
	require.NoError(t, err)
	requireBackstoryString(t, longBackstory.PersonalDescription, longDescription)
}

func TestBackstoriesTableNilUpdateDoesNotOverwriteExistingDescription(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Original description"),
	})
	require.NoError(t, err)

	updatedBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	require.Equal(t, createdBackstory.ID, updatedBackstory.ID)
	requireBackstoryString(t, updatedBackstory.PersonalDescription, "Original description")
}

func TestBackstoriesTableRequiresCharacterOwnerForUpsertGetAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              otherUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Other user should not write this"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	createdBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              owner.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Owner description"),
	})
	require.NoError(t, err)

	_, err = subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedBackstory, err := subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdBackstory.ID, deletedBackstory.ID)
}

func TestBackstoriesTableUnauthorizedUpsertDoesNotMutateExistingBackstory(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	createdBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              owner.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Original owner description"),
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              otherUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Unauthorized replacement"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedBackstory, err := subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdBackstory.ID, fetchedBackstory.ID)
	requireBackstoryString(t, fetchedBackstory.PersonalDescription, "Original owner description")
}

func TestBackstoriesTableReturnsNoRowsBeforeUpsert(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoriesTableKeepsBackstoriesScopedToRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         firstCharacter.ID,
		PersonalDescription: backstoryString("First description"),
	})
	require.NoError(t, err)

	secondBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         secondCharacter.ID,
		PersonalDescription: backstoryString("Second description"),
	})
	require.NoError(t, err)

	fetchedFirstBackstory, err := subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstBackstory.ID, fetchedFirstBackstory.ID)
	requireBackstoryString(t, fetchedFirstBackstory.PersonalDescription, "First description")

	fetchedSecondBackstory, err := subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Equal(t, secondBackstory.ID, fetchedSecondBackstory.ID)
	requireBackstoryString(t, fetchedSecondBackstory.PersonalDescription, "Second description")
}

func TestBackstoriesTableReturnsNoRowsForMissingCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	missingCharacterID := characterTestUUID("91919191-9191-9191-9191-919191919191")

	_, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         missingCharacterID,
		PersonalDescription: backstoryString("Missing character"),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoryItemsTableCreateListGetUpdateAndDeleteItem(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	backstory := createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	firstItem, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "injuries_scars",
		Title:       "Scar",
		Text:        "A scar from an expedition.",
	})
	require.NoError(t, err)
	require.True(t, firstItem.ID.Valid)
	require.Equal(t, backstory.ID, firstItem.BackstoryID)
	require.Equal(t, "injuries_scars", firstItem.Section)
	require.Equal(t, "Scar", firstItem.Title)
	require.Equal(t, "A scar from an expedition.", firstItem.Text)

	time.Sleep(5 * time.Millisecond)

	secondItem, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Careful",
		Text:        "Always checks the door twice.",
	})
	require.NoError(t, err)

	items, err := subject.queries.GetBackstoryItems(context.Background(), db.GetBackstoryItemsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, firstItem.ID, items[0].ID)
	require.Equal(t, secondItem.ID, items[1].ID)

	fetchedItem, err := subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: firstItem.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstItem.ID, fetchedItem.ID)

	updatedItem, err := subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		Section:         "encounters",
		Title:           "Deep One",
		Text:            "Saw something in Innsmouth.",
		BackstoryItemID: firstItem.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstItem.ID, updatedItem.ID)
	require.Equal(t, "encounters", updatedItem.Section)
	require.Equal(t, "Deep One", updatedItem.Title)
	require.Equal(t, "Saw something in Innsmouth.", updatedItem.Text)

	deletedItem, err := subject.queries.DeleteBackstoryItem(context.Background(), db.DeleteBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: firstItem.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstItem.ID, deletedItem.ID)

	_, err = subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: firstItem.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoryItemsTableListReturnsEmptyWhenBackstoryHasNoItems(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	items, err := subject.queries.GetBackstoryItems(context.Background(), db.GetBackstoryItemsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestBackstoryItemsTableRequiresBackstoryBeforeCreateOrList(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "No backstory",
		Text:        "This should not be created.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	items, err := subject.queries.GetBackstoryItems(context.Background(), db.GetBackstoryItemsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestBackstoryItemsTableRequiresCharacterOwnerForCreateListGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)
	createBackstoryTestBackstory(t, subject, owner.ID, character.ID)

	_, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Hidden",
		Text:        "No access.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Visible",
		Text:        "Owner access.",
	})
	require.NoError(t, err)

	items, err := subject.queries.GetBackstoryItems(context.Background(), db.GetBackstoryItemsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, items)

	_, err = subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          otherUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          otherUser.ID,
		CharacterID:     character.ID,
		Section:         "encounters",
		Title:           "Changed",
		Text:            "Should not change.",
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstoryItem(context.Background(), db.DeleteBackstoryItemParams{
		UserID:          otherUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoryItemsTableRequiresMatchingCharacterForGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, firstCharacter.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, secondCharacter.ID)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
		Section:     "traits",
		Title:       "First character item",
		Text:        "Scoped to first character.",
	})
	require.NoError(t, err)

	_, err = subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     secondCharacter.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     secondCharacter.ID,
		Section:         "encounters",
		Title:           "Wrong character",
		Text:            "Should not update.",
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstoryItem(context.Background(), db.DeleteBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     secondCharacter.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoryItemsTableRejectsInvalidSectionOnCreateAndUpdate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	_, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "invalid_section",
		Title:       "Invalid",
		Text:        "Invalid section.",
	})
	requirePostgresErrorCode(t, err, "23514")

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Valid",
		Text:        "Valid section.",
	})
	require.NoError(t, err)

	_, err = subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		Section:         "invalid_section",
		Title:           "Invalid",
		Text:            "Invalid section.",
		BackstoryItemID: item.ID,
	})
	requirePostgresErrorCode(t, err, "23514")
}

func TestBackstoryItemsTableAllowsEveryValidSection(t *testing.T) {
	validSections := []string{
		"injuries_scars",
		"phobias_manias",
		"arcane_tomes_spells",
		"encounters",
		"ideology_beliefs",
		"significant_people",
		"meaningful_locations",
		"treasured_possessions",
		"traits",
	}

	for _, section := range validSections {
		t.Run(section, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)
			createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

			item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
				UserID:      testUser.ID,
				CharacterID: character.ID,
				Section:     section,
				Title:       section,
				Text:        "Valid section.",
			})
			require.NoError(t, err)
			require.Equal(t, section, item.Section)
		})
	}
}

func TestBackstoryItemsTableRejectsTooLongSectionAndTitle(t *testing.T) {
	tests := []struct {
		name   string
		params func(userID string, characterID pgtype.UUID) db.CreateBackstoryItemParams
	}{
		{
			name: "section",
			params: func(userID string, characterID pgtype.UUID) db.CreateBackstoryItemParams {
				return db.CreateBackstoryItemParams{
					UserID:      userID,
					CharacterID: characterID,
					Section:     strings.Repeat("a", 33),
					Title:       "Title",
					Text:        "Text",
				}
			},
		},
		{
			name: "title",
			params: func(userID string, characterID pgtype.UUID) db.CreateBackstoryItemParams {
				return db.CreateBackstoryItemParams{
					UserID:      userID,
					CharacterID: characterID,
					Section:     "traits",
					Title:       strings.Repeat("a", 256),
					Text:        "Text",
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)
			createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

			_, err := subject.queries.CreateBackstoryItem(context.Background(), tc.params(testUser.ID, character.ID))
			requirePostgresErrorCode(t, err, "22001")
		})
	}
}

func TestBackstoryItemsTableAllowsTitleAtLengthLimit(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)
	title := strings.Repeat("a", 255)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       title,
		Text:        "Boundary title.",
	})
	require.NoError(t, err)
	require.Equal(t, title, item.Title)
}

func TestBackstoryItemsTableAllowsEmptyTitleAndTextAndLongText(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)
	longText := strings.Repeat("text ", 1000)

	emptyItem, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "",
		Text:        "",
	})
	require.NoError(t, err)
	require.Equal(t, "", emptyItem.Title)
	require.Equal(t, "", emptyItem.Text)

	longItem, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Long text",
		Text:        longText,
	})
	require.NoError(t, err)
	require.Equal(t, longText, longItem.Text)
}

func TestBackstoryItemsTableReturnsNoRowsForMissingItem(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)
	missingItemID := characterTestUUID("92929292-9292-9292-9292-929292929292")

	_, err := subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: missingItemID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		Section:         "traits",
		Title:           "Missing",
		Text:            "Missing item.",
		BackstoryItemID: missingItemID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteBackstoryItem(context.Background(), db.DeleteBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: missingItemID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoryItemsTableFailedUpdateDoesNotMutateExistingItem(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)
	createBackstoryTestBackstory(t, subject, owner.ID, character.ID)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Original title",
		Text:        "Original text.",
	})
	require.NoError(t, err)

	_, err = subject.queries.UpdateBackstoryItem(context.Background(), db.UpdateBackstoryItemParams{
		UserID:          otherUser.ID,
		CharacterID:     character.ID,
		Section:         "encounters",
		Title:           "Unauthorized title",
		Text:            "Unauthorized text.",
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedItem, err := subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          owner.ID,
		CharacterID:     character.ID,
		BackstoryItemID: item.ID,
	})
	require.NoError(t, err)
	require.Equal(t, "traits", fetchedItem.Section)
	require.Equal(t, "Original title", fetchedItem.Title)
	require.Equal(t, "Original text.", fetchedItem.Text)
}

func TestBackstoryItemsTableDeletingBackstoryCascadesItems(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Cascade",
		Text:        "Deleted with backstory.",
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestBackstoriesTableDeleteReturnsDeletedValuesAndAllowsRecreate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdBackstory := createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	deletedBackstory, err := subject.queries.DeleteBackstory(context.Background(), db.DeleteBackstoryParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdBackstory.ID, deletedBackstory.ID)
	requireBackstoryString(t, deletedBackstory.PersonalDescription, "Test backstory")

	recreatedBackstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              testUser.ID,
		CharacterID:         character.ID,
		PersonalDescription: backstoryString("Recreated backstory"),
	})
	require.NoError(t, err)
	require.NotEqual(t, deletedBackstory.ID, recreatedBackstory.ID)
	requireBackstoryString(t, recreatedBackstory.PersonalDescription, "Recreated backstory")
}

func TestBackstoriesTableDeletingCharacterCascadesBackstoryAndItems(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	createBackstoryTestBackstory(t, subject, testUser.ID, character.ID)

	item, err := subject.queries.CreateBackstoryItem(context.Background(), db.CreateBackstoryItemParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Section:     "traits",
		Title:       "Cascade",
		Text:        "Deleted with character.",
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetBackstoryByCharacter(context.Background(), db.GetBackstoryByCharacterParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetBackstoryItem(context.Background(), db.GetBackstoryItemParams{
		UserID:          testUser.ID,
		CharacterID:     character.ID,
		BackstoryItemID: item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func createBackstoryTestBackstory(t *testing.T, subject *characterIntegrationSubject, userID string, characterID pgtype.UUID) db.Backstory {
	t.Helper()

	backstory, err := subject.queries.UpsertBackstory(context.Background(), db.UpsertBackstoryParams{
		UserID:              userID,
		CharacterID:         characterID,
		PersonalDescription: backstoryString("Test backstory"),
	})
	require.NoError(t, err)

	return backstory
}

func backstoryString(value string) *string {
	return &value
}

func requireBackstoryString(t *testing.T, actual *string, expected string) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}
