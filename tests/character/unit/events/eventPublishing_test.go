package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
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
				_, err := service.GetCharacter(ctx, characterDTO.GetCharacterInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "create character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateCharacter(ctx, characterDTO.CreateCharacterInput{UserID: testUserID})
				return err
			},
			expectedEvent: characterEvents.CharacterCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "update character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateCharacter(ctx, characterDTO.UpdateCharacterInput{UserID: testUserID, ID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, Name: "Dr. Armitage"},
		},
		{
			name: "replace character portrait",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.ReplacePortrait(ctx, characterDTO.ReplacePortraitInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterPortraitReplaceSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacter(ctx, characterDTO.DeleteCharacterInput{UserID: testUserID, ID: characterID})
			},
			expectedEvent: characterEvents.CharacterDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetHealth(ctx, healthDTO.GetHealthInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertHealth(ctx, healthDTO.UpsertHealthInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteHealth(ctx, healthDTO.DeleteHealthInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterHealthDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSanity(ctx, sanityDTO.GetSanityInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertSanity(ctx, sanityDTO.UpsertSanityInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSanity(ctx, sanityDTO.DeleteSanityInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterSanityDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetMagic(ctx, magicDTO.GetMagicInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertMagic(ctx, magicDTO.UpsertMagicInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteMagic(ctx, magicDTO.DeleteMagicInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterMagicDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetLuck(ctx, luckDTO.GetLuckInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertLuck(ctx, luckDTO.UpsertLuckInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteLuck(ctx, luckDTO.DeleteLuckInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterLuckDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetFinances(ctx, financesDTO.GetFinancesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertFinances(ctx, financesDTO.UpsertFinancesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteFinances(ctx, financesDTO.DeleteFinancesInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterFinancesDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstory(ctx, backstoriesDTO.GetBackstoryInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertBackstory(ctx, backstoriesDTO.UpsertBackstoryInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstory(ctx, backstoriesDTO.DeleteBackstoryInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterBackstoryDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "list backstory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItems(ctx, backstoriesDTO.GetBackstoryItemsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemsListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItem(ctx, backstoriesDTO.GetBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "create backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateBackstoryItem(ctx, backstoriesDTO.CreateBackstoryItemInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "update backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateBackstoryItem(ctx, backstoriesDTO.UpdateBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Section: "ideology_beliefs", Title: "Old Motto"},
		},
		{
			name: "delete backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstoryItem(ctx, backstoriesDTO.DeleteBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterBackstoryItemDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID},
		},
		{
			name: "list skills",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkills(ctx, skillsDTO.GetSkillsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillsListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkill(ctx, skillsDTO.GetSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "create skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateSkill(ctx, skillsDTO.CreateSkillInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "update skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateSkill(ctx, skillsDTO.UpdateSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Name: "Library Use"},
		},
		{
			name: "delete skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSkill(ctx, skillsDTO.DeleteSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
			},
			expectedEvent: characterEvents.CharacterSkillDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID},
		},
		{
			name: "get derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetDerivedStats(ctx, derivedStatsDTO.GetDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertDerivedStats(ctx, derivedStatsDTO.UpsertDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteDerivedStats(ctx, derivedStatsDTO.DeleteDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterDerivedStatsDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "get characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacteristics(ctx, characteristicsDTO.GetCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsGetSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "upsert characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertCharacteristics(ctx, characteristicsDTO.UpsertCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsUpsertSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "delete characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacteristics(ctx, characteristicsDTO.DeleteCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterCharacteristicsDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID},
		},
		{
			name: "list notes",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNotes(ctx, notesDTO.GetNotesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNotesListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 2},
		},
		{
			name: "get note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNote(ctx, notesDTO.GetNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "create note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateNote(ctx, notesDTO.CreateNoteInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "update note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateNote(ctx, notesDTO.UpdateNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Title: "Session Note"},
		},
		{
			name: "delete note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteNote(ctx, notesDTO.DeleteNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
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
				_, err := service.GetCharacter(ctx, characterDTO.GetCharacterInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "create character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateCharacter(ctx, characterDTO.CreateCharacterInput{UserID: testUserID})
				return err
			},
			expectedEvent: characterEvents.CharacterCreateFailed{UserID: testUserID, Err: expectedErr},
		},
		{
			name: "update character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateCharacter(ctx, characterDTO.UpdateCharacterInput{UserID: testUserID, ID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "replace character portrait",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.ReplacePortrait(ctx, characterDTO.ReplacePortraitInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterPortraitReplaceFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete character",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacter(ctx, characterDTO.DeleteCharacterInput{UserID: testUserID, ID: characterID})
			},
			expectedEvent: characterEvents.CharacterDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetHealth(ctx, healthDTO.GetHealthInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertHealth(ctx, healthDTO.UpsertHealthInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterHealthUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete health",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteHealth(ctx, healthDTO.DeleteHealthInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterHealthDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSanity(ctx, sanityDTO.GetSanityInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertSanity(ctx, sanityDTO.UpsertSanityInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSanityUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete sanity",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSanity(ctx, sanityDTO.DeleteSanityInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterSanityDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetMagic(ctx, magicDTO.GetMagicInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertMagic(ctx, magicDTO.UpsertMagicInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterMagicUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete magic",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteMagic(ctx, magicDTO.DeleteMagicInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterMagicDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetLuck(ctx, luckDTO.GetLuckInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertLuck(ctx, luckDTO.UpsertLuckInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterLuckUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete luck",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteLuck(ctx, luckDTO.DeleteLuckInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterLuckDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetFinances(ctx, financesDTO.GetFinancesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertFinances(ctx, financesDTO.UpsertFinancesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterFinancesUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete finances",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteFinances(ctx, financesDTO.DeleteFinancesInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterFinancesDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstory(ctx, backstoriesDTO.GetBackstoryInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertBackstory(ctx, backstoriesDTO.UpsertBackstoryInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete backstory",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstory(ctx, backstoriesDTO.DeleteBackstoryInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterBackstoryDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "list backstory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItems(ctx, backstoriesDTO.GetBackstoryItemsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemsListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetBackstoryItem(ctx, backstoriesDTO.GetBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemGetFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "create backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateBackstoryItem(ctx, backstoriesDTO.CreateBackstoryItemInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateBackstoryItem(ctx, backstoriesDTO.UpdateBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterBackstoryItemUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "delete backstory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteBackstoryItem(ctx, backstoriesDTO.DeleteBackstoryItemInput{UserID: testUserID, CharacterID: characterID, BackstoryItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterBackstoryItemDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, BackstoryItemID: testItemID, Err: expectedErr},
		},
		{
			name: "list skills",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkills(ctx, skillsDTO.GetSkillsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillsListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetSkill(ctx, skillsDTO.GetSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillGetFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "create skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateSkill(ctx, skillsDTO.CreateSkillInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateSkill(ctx, skillsDTO.UpdateSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
				return err
			},
			expectedEvent: characterEvents.CharacterSkillUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "delete skill",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteSkill(ctx, skillsDTO.DeleteSkillInput{UserID: testUserID, CharacterID: characterID, SkillID: skillID})
			},
			expectedEvent: characterEvents.CharacterSkillDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, SkillID: testSkillID, Err: expectedErr},
		},
		{
			name: "get derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetDerivedStats(ctx, derivedStatsDTO.GetDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertDerivedStats(ctx, derivedStatsDTO.UpsertDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterDerivedStatsUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete derived stats",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteDerivedStats(ctx, derivedStatsDTO.DeleteDerivedStatsInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterDerivedStatsDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetCharacteristics(ctx, characteristicsDTO.GetCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsGetFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "upsert characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpsertCharacteristics(ctx, characteristicsDTO.UpsertCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterCharacteristicsUpsertFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "delete characteristics",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteCharacteristics(ctx, characteristicsDTO.DeleteCharacteristicsInput{UserID: testUserID, CharacterID: characterID})
			},
			expectedEvent: characterEvents.CharacterCharacteristicsDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "list notes",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNotes(ctx, notesDTO.GetNotesInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNotesListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetNote(ctx, notesDTO.GetNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteGetFailed{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Err: expectedErr},
		},
		{
			name: "create note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateNote(ctx, notesDTO.CreateNoteInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateNote(ctx, notesDTO.UpdateNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
				return err
			},
			expectedEvent: characterEvents.CharacterNoteUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, NoteID: testNoteID, Err: expectedErr},
		},
		{
			name: "delete note",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteNote(ctx, notesDTO.DeleteNoteInput{UserID: testUserID, CharacterID: characterID, NoteID: noteID})
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
