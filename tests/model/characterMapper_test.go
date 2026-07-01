package tests

import (
	"testing"

	characterModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestToShortCharacterModelCopiesAllCharacterFields(t *testing.T) {
	character := testCharacter()

	result := characterModel.ToCharacterShortModel(character)

	requireSameShortCharacter(t, character, result)
}

func TestToShortCharacterModelPreservesNilOptionalFields(t *testing.T) {
	character := testCharacter()
	character.PlayerName = nil
	character.Occupation = nil
	character.Age = nil
	character.Sex = nil
	character.Residence = nil
	character.Birthplace = nil

	result := characterModel.ToCharacterShortModel(character)

	require.Nil(t, result.PlayerName)
	require.Nil(t, result.Occupation)
	require.Nil(t, result.Age)
	require.Nil(t, result.Sex)
	require.Nil(t, result.Residence)
	require.Nil(t, result.Birthplace)
}

func TestToFullCharacterModelLeavesOptionalSectionsEmptyWhenIDsAreInvalid(t *testing.T) {
	character := testCharacter()

	result := characterModel.ToCharacterModel(characterModel.CharacterDBData{
		Character:       character,
		Characteristics: db.Characteristic{ID: invalidUUID()},
		DerivedStats:    db.DerivedStat{ID: invalidUUID()},
		HP:              db.HealthState{ID: invalidUUID()},
		MP:              db.MagicState{ID: invalidUUID()},
		Sanity:          db.SanityState{ID: invalidUUID()},
		Luck:            db.LuckState{ID: invalidUUID()},
	})

	requireSameShortCharacter(t, character, result.CharacterShortModel)
	require.False(t, result.Characteristics.ID.Valid)
	require.False(t, result.DerivedStats.ID.Valid)
	require.False(t, result.HP.ID.Valid)
	require.False(t, result.MP.ID.Valid)
	require.False(t, result.Sanity.ID.Valid)
	require.False(t, result.Luck.ID.Valid)
	require.Empty(t, result.Skills)
	require.Empty(t, result.Notes)
	require.False(t, result.Backstory.ID.Valid)
	require.False(t, result.Finances.ID.Valid)
}

func TestToFullCharacterModelMapsAllPresentSections(t *testing.T) {
	character := testCharacter()
	skill := testSkillRow()
	creditSkill := testSpecializedSkillRow()
	backstory := testBackstory()
	item := testBackstoryItem()
	finance := testFinance()
	finance.CreditRatingSkillID = creditSkill.ID
	note := db.Note{
		ID:          testUUID("99999999-9999-9999-9999-999999999999"),
		CharacterID: character.ID,
		Title:       "Session note",
		Body:        "Found a hidden index.",
		CreatedAt:   testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:   testTimestamptz("2026-06-07 13:00:00+03"),
	}
	characteristics := db.Characteristic{ID: testUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), CharacterID: character.ID}
	derivedStats := db.DerivedStat{ID: testUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), CharacterID: character.ID}
	hp := db.HealthState{ID: testUUID("cccccccc-cccc-cccc-cccc-cccccccccccc"), CharacterID: character.ID, MaxHp: 10, CurrentHp: 7}
	mp := db.MagicState{ID: testUUID("dddddddd-dddd-dddd-dddd-dddddddddddd"), CharacterID: character.ID, MaxMp: 10, CurrentMp: 5}
	sanity := db.SanityState{ID: testUUID("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"), CharacterID: character.ID, MaxSanity: 60, CurrentSanity: 40}
	luck := db.LuckState{ID: testUUID("ffffffff-ffff-ffff-ffff-ffffffffffff"), CharacterID: character.ID, StartingLuck: 50, CurrentLuck: 30}

	result := characterModel.ToCharacterModel(characterModel.CharacterDBData{
		Character:       character,
		Characteristics: characteristics,
		DerivedStats:    derivedStats,
		HP:              hp,
		MP:              mp,
		Sanity:          sanity,
		Luck:            luck,
		Skills:          []db.GetSkillsRow{skill, creditSkill},
		Backstory:       &backstory,
		BackstoryItems:  []db.BackstoryItem{item},
		Finances:        &finance,
		Notes:           []db.Note{note},
	})

	requireSameShortCharacter(t, character, result.CharacterShortModel)
	require.Equal(t, characteristics, result.Characteristics)
	require.Equal(t, derivedStats, result.DerivedStats)
	require.Equal(t, hp, result.HP)
	require.Equal(t, mp, result.MP)
	require.Equal(t, sanity, result.Sanity)
	require.Equal(t, luck, result.Luck)
	require.Len(t, result.Skills, 2)
	requireSameSkill(t, skill, result.Skills[0])
	requireSameSkill(t, creditSkill, result.Skills[1])
	require.Equal(t, backstory.ID, result.Backstory.ID)
	require.Len(t, result.Backstory.Items, 1)
	require.Equal(t, item.ID, result.Backstory.Items[0].ID)
	require.Equal(t, finance.ID, result.Finances.ID)
	require.NotNil(t, result.Finances.CreditRating)
	require.Equal(t, creditSkill.ID, result.Finances.CreditRating.ID)
	require.Equal(t, []db.Note{note}, result.Notes)
}

func TestToFullCharacterModelLeavesCreditRatingNilWhenFinanceSkillDoesNotMatch(t *testing.T) {
	finance := testFinance()
	finance.CreditRatingSkillID = testUUID("abababab-abab-abab-abab-abababababab")

	result := characterModel.ToCharacterModel(characterModel.CharacterDBData{
		Character: testCharacter(),
		Skills:    []db.GetSkillsRow{testSkillRow()},
		Finances:  &finance,
	})

	require.Equal(t, finance.ID, result.Finances.ID)
	require.Nil(t, result.Finances.CreditRating)
}
