package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/stretchr/testify/require"
)

type eventPublishingCase struct {
	name          string
	call          func(context.Context, *characterServices.EventPublishingCharacterService) error
	expectedEvent events.Event
}

func TestEventPublishingCharacterServicePublishesSuccessEvents(t *testing.T) {
	characterID := testUUID(testCharacterID)
	itemID := testUUID(testItemID)
	skillID := testUUID(testSkillID)
	noteID := testUUID(testNoteID)

	cases := []eventPublishingCase{
		{
			name: "get all characters",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetAllCharacters(ctx, testUserID)
				return err
			},
			expectedEvent: characterEvents.CharactersListSucceeded{UserID: testUserID, Count: 2},
		},
		{
			name: "get character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacter(ctx, model.GetCharacterInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "create character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateCharacter(ctx, db.CreateCharacterParams{UserID: testUserID})
				return err
			},
			expectedEvent: characterEvents.CharacterCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "update character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateCharacter(ctx, db.UpdateCharacterParams{UserID: testUserID, ID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "delete character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacter(ctx, db.DeleteCharacterParams{UserID: testUserID, ID: characterID})
			},
			expectedEvent: characterEvents.CharacterDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetHealth(ctx, db.GetHealthStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertHealth(ctx, db.UpsertHealthStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteHealth(ctx, db.DeleteHealthStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterHealthDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSanity(ctx, db.GetSanityStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertSanity(ctx, db.UpsertSanityStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSanity(ctx, db.DeleteSanityStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterSanityDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetMagic(ctx, db.GetMagicStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertMagic(ctx, db.UpsertMagicStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteMagic(ctx, db.DeleteMagicStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterMagicDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetLuck(ctx, db.GetLuckStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertLuck(ctx, db.UpsertLuckStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteLuck(ctx, db.DeleteLuckStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterLuckDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetFinances(ctx, db.GetFinancesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertFinances(ctx, db.UpsertFinancesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteFinances(ctx, db.DeleteFinancesParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterFinancesDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstory(ctx, db.GetBackstoryByCharacterParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertBackstory(ctx, db.UpsertBackstoryParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstory(ctx, db.DeleteBackstoryParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterBackstoryDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "list backstory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemsListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItem(ctx, db.GetBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "create backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateBackstoryItem(ctx, db.CreateBackstoryItemParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "update backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateBackstoryItem(ctx, db.UpdateBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "delete backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstoryItem(ctx, db.DeleteBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterBackstoryItemDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID},
		},
		{
			name: "list skills",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkills(ctx, db.GetCharacterSkillsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillsListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkill(ctx, db.GetCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "create skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateSkill(ctx, db.CreateCharacterSkillParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "update skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateSkill(ctx, db.UpdateCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "delete skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSkill(ctx, db.DeleteCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
			},
			expectedEvent: characterEvents.CharacterSkillDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID},
		},
		{
			name: "get derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetDerivedStats(ctx, db.GetDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertDerivedStats(ctx, db.UpsertDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteDerivedStats(ctx, db.DeleteDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterDerivedStatsDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacteristics(ctx, db.GetCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertCharacteristics(ctx, db.UpsertCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacteristics(ctx, db.DeleteCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterCharacteristicsDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "list notes",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNotes(ctx, db.GetNotesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNotesListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNote(ctx, db.GetNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "create note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateNote(ctx, db.CreateNoteParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "update note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateNote(ctx, db.UpdateNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "delete note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteNote(ctx, db.DeleteNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
			},
			expectedEvent: characterEvents.CharacterNoteDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, publisher, service := newEventPublishingTestSubject()

			err := tc.call(context.Background(), service)

			require.NoError(t, err)
			requirePublishedEvent(t, publisher, tc.expectedEvent)
		})
	}
}

func TestEventPublishingCharacterServicePublishesFailureEvents(t *testing.T) {
	characterID := testUUID(testCharacterID)
	itemID := testUUID(testItemID)
	skillID := testUUID(testSkillID)
	noteID := testUUID(testNoteID)
	expectedErr := errors.New("base service failed")

	cases := []eventPublishingCase{
		{
			name: "get all characters",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetAllCharacters(ctx, testUserID)
				return err
			},
			expectedEvent: characterEvents.CharactersListFailed{UserID: testUserID, Err: expectedErr},
		},
		{
			name: "get character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacter(ctx, model.GetCharacterInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "create character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateCharacter(ctx, db.CreateCharacterParams{UserID: testUserID})
				return err
			},
			expectedEvent: characterEvents.CharacterCreateFailed{UserID: testUserID, Err: expectedErr},
		},
		{
			name: "update character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateCharacter(ctx, db.UpdateCharacterParams{UserID: testUserID, ID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacter(ctx, db.DeleteCharacterParams{UserID: testUserID, ID: characterID})
			},
			expectedEvent: characterEvents.CharacterDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetHealth(ctx, db.GetHealthStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertHealth(ctx, db.UpsertHealthStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteHealth(ctx, db.DeleteHealthStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterHealthDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSanity(ctx, db.GetSanityStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertSanity(ctx, db.UpsertSanityStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSanity(ctx, db.DeleteSanityStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterSanityDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetMagic(ctx, db.GetMagicStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertMagic(ctx, db.UpsertMagicStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteMagic(ctx, db.DeleteMagicStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterMagicDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetLuck(ctx, db.GetLuckStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertLuck(ctx, db.UpsertLuckStateParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteLuck(ctx, db.DeleteLuckStateParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterLuckDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetFinances(ctx, db.GetFinancesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertFinances(ctx, db.UpsertFinancesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteFinances(ctx, db.DeleteFinancesParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterFinancesDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstory(ctx, db.GetBackstoryByCharacterParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertBackstory(ctx, db.UpsertBackstoryParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstory(ctx, db.DeleteBackstoryParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterBackstoryDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "list backstory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemsListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItem(ctx, db.GetBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemGetFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "create backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateBackstoryItem(ctx, db.CreateBackstoryItemParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateBackstoryItem(ctx, db.UpdateBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "delete backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstoryItem(ctx, db.DeleteBackstoryItemParams{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterBackstoryItemDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "list skills",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkills(ctx, db.GetCharacterSkillsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillsListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkill(ctx, db.GetCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillGetFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "create skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateSkill(ctx, db.CreateCharacterSkillParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateSkill(ctx, db.UpdateCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "delete skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSkill(ctx, db.DeleteCharacterSkillParams{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
			},
			expectedEvent: characterEvents.CharacterSkillDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "get derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetDerivedStats(ctx, db.GetDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertDerivedStats(ctx, db.UpsertDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteDerivedStats(ctx, db.DeleteDerivedStatsParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterDerivedStatsDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacteristics(ctx, db.GetCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertCharacteristics(ctx, db.UpsertCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacteristics(ctx, db.DeleteCharacteristicsParams{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterCharacteristicsDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "list notes",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNotes(ctx, db.GetNotesParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNotesListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNote(ctx, db.GetNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteGetFailed{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Err: expectedErr},
		},
		{
			name: "create note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateNote(ctx, db.CreateNoteParams{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateNote(ctx, db.UpdateNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Err: expectedErr},
		},
		{
			name: "delete note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteNote(ctx, db.DeleteNoteParams{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
			},
			expectedEvent: characterEvents.CharacterNoteDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Err: expectedErr},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next, publisher, service := newEventPublishingTestSubject()
			next.Err = expectedErr

			err := tc.call(context.Background(), service)

			require.ErrorIs(t, err, expectedErr)
			requirePublishedEvent(t, publisher, tc.expectedEvent)
		})
	}
}
