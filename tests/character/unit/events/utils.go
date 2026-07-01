package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

const (
	testUserID      = "user_1"
	testCharacterID = "11111111-1111-1111-1111-111111111111"
	testItemID      = "22222222-2222-2222-2222-222222222222"
	testSkillID     = "33333333-3333-3333-3333-333333333333"
	testNoteID      = "44444444-4444-4444-4444-444444444444"
)

func newEventPublishingTestSubject() (*FakeCharacterService, *FakeEventPublisher, *characterServices.EventPublishingCharacterService) {
	next := newFakeCharacterService()
	publisher := &FakeEventPublisher{}

	return next, publisher, characterServices.NewEventPublishingCharacterService(next, publisher)
}

func newFakeCharacterService() *FakeCharacterService {
	return &FakeCharacterService{
		Characters: []characterDTO.CharacterShortModel{
			testCharacterShortModel(),
			{ID: testUUID("55555555-5555-5555-5555-555555555555"), UserID: testUserID, Name: "Second Character"},
		},
		Character:       testCharacterModel(),
		Backstory:       backstoriesDTO.BackstoryModel{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		BackstoryItems:  []backstoriesDTO.BackstoryItemModel{testBackstoryItemModel(), {ID: testUUID("66666666-6666-6666-6666-666666666666")}},
		BackstoryItem:   testBackstoryItemModel(),
		Skills:          []skillsDTO.SkillModel{testSkillModel(), {ID: testUUID("77777777-7777-7777-7777-777777777777"), Name: "Spot Hidden"}},
		Skill:           testSkillModel(),
		Notes:           []db.Note{testNote(), {ID: testUUID("88888888-8888-8888-8888-888888888888"), Title: "Second Note"}},
		Note:            testNote(),
		Health:          db.HealthState{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		Sanity:          db.SanityState{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		Magic:           db.MagicState{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		Luck:            db.LuckState{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		Finances:        db.Finance{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		DerivedStats:    db.DerivedStat{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
		Characteristics: db.Characteristic{ID: testUUID(testItemID), CharacterID: testUUID(testCharacterID)},
	}
}

func testUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func testCharacterShortModel() characterDTO.CharacterShortModel {
	return characterDTO.CharacterShortModel{
		ID:     testUUID(testCharacterID),
		UserID: testUserID,
		Name:   "Dr. Armitage",
	}
}

func testCharacterModel() characterDTO.CharacterModel {
	return characterDTO.CharacterModel{
		CharacterShortModel: testCharacterShortModel(),
	}
}

func testBackstoryItemModel() backstoriesDTO.BackstoryItemModel {
	return backstoriesDTO.BackstoryItemModel{
		ID:      testUUID(testItemID),
		Section: "ideology_beliefs",
		Title:   "Old Motto",
		Text:    "Knowledge has a price.",
	}
}

func testSkillModel() skillsDTO.SkillModel {
	return skillsDTO.SkillModel{
		ID:   testUUID(testSkillID),
		Name: "Library Use",
	}
}

func testNote() db.Note {
	return db.Note{
		ID:          testUUID(testNoteID),
		CharacterID: testUUID(testCharacterID),
		Title:       "Session Note",
		Body:        "Found the hidden index.",
	}
}

func requirePublishedEvent(t *testing.T, publisher *FakeEventPublisher, expected events.Event) {
	t.Helper()

	require.Len(t, publisher.Events, 1)
	require.Equal(t, expected, publisher.Events[0])
}
