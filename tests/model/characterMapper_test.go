package tests

import (
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestToShortCharacterModelCopiesAllCharacterFields(t *testing.T) {
	character := testCharacter()

	result := characterDTO.ToCharacterShortModel(character)

	requireSameShortCharacter(t, character, result)
}

func TestToShortCharacterModelPreservesNilOptionalFields(t *testing.T) {
	character := testCharacter()
	character.Occupation = nil
	character.Age = nil
	character.Sex = nil
	character.Residence = nil
	character.Birthplace = nil
	character.PortraitKey = nil

	result := characterDTO.ToCharacterShortModel(character)

	require.Nil(t, result.Occupation)
	require.Nil(t, result.Age)
	require.Nil(t, result.Sex)
	require.Nil(t, result.Residence)
	require.Nil(t, result.Birthplace)
	require.Nil(t, result.PortraitUrl)
}

func TestToCharacterSummaryModelMapsBaseFieldsAndStats(t *testing.T) {
	portraitKey := "portraits/11111111-1111-1111-1111-111111111111.webp"
	row := db.GetAllUserCharacterCardsRow{
		ID:            testUUID("11111111-1111-1111-1111-111111111111"),
		Name:          "Dr. Armitage",
		Occupation:    strPtr("Antiquarian"),
		Age:           int16Ptr(42),
		Sex:           strPtr(""),
		Residence:     strPtr("Arkham"),
		PortraitKey:   &portraitKey,
		CurrentHp:     7,
		MaxHp:         12,
		CurrentMp:     4,
		MaxMp:         9,
		CurrentSanity: 33,
		MaxSanity:     60,
		CurrentLuck:   20,
		StartingLuck:  45,
	}

	result := characterDTO.ToCharacterSummaryModel(row)

	require.Equal(t, row.ID, result.ID)
	require.Equal(t, row.Name, result.Name)
	require.Equal(t, row.Occupation, result.Occupation)
	require.Equal(t, row.Age, result.Age)
	require.Equal(t, row.Sex, result.Sex)
	require.Equal(t, row.Residence, result.Residence)
	require.Nil(t, result.PortraitUrl)
	require.Equal(t, int16(7), result.HP.Current)
	require.Equal(t, int16(12), result.HP.Max)
	require.Equal(t, int16(4), result.MP.Current)
	require.Equal(t, int16(9), result.MP.Max)
	require.Equal(t, int16(33), result.Sanity.Current)
	require.Equal(t, int16(60), result.Sanity.Max)
	require.Equal(t, int16(20), result.Luck.Current)
	require.Equal(t, int16(45), result.Luck.Starting)
}

func TestToFullCharacterModelLeavesOptionalSectionsEmptyWhenIDsAreInvalid(t *testing.T) {
	character := testCharacter()

	result := characterDTO.ToCharacterModel(characterDTO.CharacterDBData{
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
	creditSkill := testSkillRow()
	creditSkill.ID = testUUID("66666666-6666-6666-6666-666666666666")
	creditSkill.Name = "Credit Rating"
	creditSkill.IsProtected = false
	creditSkill.BaseRule = nil
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

	result := characterDTO.ToCharacterModel(characterDTO.CharacterDBData{
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
	requireEqualCharacteristic(t, characteristics, result.Characteristics)
	requireEqualDerivedStats(t, derivedStats, result.DerivedStats)
	requireEqualHealth(t, hp, result.HP)
	requireEqualMagic(t, mp, result.MP)
	requireEqualSanity(t, sanity, result.Sanity)
	requireEqualLuck(t, luck, result.Luck)
	require.Len(t, result.Skills, 2)
	requireSameSkill(t, skill, result.Skills[0])
	requireSameSkill(t, creditSkill, result.Skills[1])
	require.Equal(t, backstory.ID, result.Backstory.ID)
	require.Len(t, result.Backstory.Items, 1)
	require.Equal(t, item.ID, result.Backstory.Items[0].ID)
	require.Equal(t, finance.ID, result.Finances.ID)
	require.NotNil(t, result.Finances.CreditRating)
	require.Equal(t, creditSkill.ID, result.Finances.CreditRating.ID)
	requireEqualNotes(t, []db.Note{note}, result.Notes)
}

func TestToFullCharacterModelLeavesCreditRatingNilWhenFinanceSkillDoesNotMatch(t *testing.T) {
	finance := testFinance()
	finance.CreditRatingSkillID = testUUID("abababab-abab-abab-abab-abababababab")

	result := characterDTO.ToCharacterModel(characterDTO.CharacterDBData{
		Character: testCharacter(),
		Skills:    []db.GetSkillsRow{testSkillRow()},
		Finances:  &finance,
	})

	require.Equal(t, finance.ID, result.Finances.ID)
	require.Nil(t, result.Finances.CreditRating)
}

func requireEqualCharacteristic(t *testing.T, expected db.Characteristic, actual characteristicsDTO.CharacteristicsModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.Strength, actual.Strength)
	require.Equal(t, expected.Constitution, actual.Constitution)
	require.Equal(t, expected.Size, actual.Size)
	require.Equal(t, expected.Dexterity, actual.Dexterity)
	require.Equal(t, expected.Appearance, actual.Appearance)
	require.Equal(t, expected.Intelligence, actual.Intelligence)
	require.Equal(t, expected.Power, actual.Power)
	require.Equal(t, expected.Education, actual.Education)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualDerivedStats(t *testing.T, expected db.DerivedStat, actual derivedStatsDTO.DerivedStatsModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.Speed, actual.Speed)
	require.Equal(t, expected.Physique, actual.Physique)
	require.Equal(t, expected.DamageBonus, actual.DamageBonus)
	require.Equal(t, expected.DodgeValue, actual.DodgeValue)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualHealth(t *testing.T, expected db.HealthState, actual healthDTO.HealthModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxHp, actual.MaxHp)
	require.Equal(t, expected.CurrentHp, actual.CurrentHp)
	require.Equal(t, expected.MajorWound, actual.MajorWound)
	require.Equal(t, expected.Unconscious, actual.Unconscious)
	require.Equal(t, expected.Dying, actual.Dying)
	require.Equal(t, expected.Dead, actual.Dead)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualMagic(t *testing.T, expected db.MagicState, actual magicDTO.MagicModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxMp, actual.MaxMp)
	require.Equal(t, expected.CurrentMp, actual.CurrentMp)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualSanity(t *testing.T, expected db.SanityState, actual sanityDTO.SanityModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.MaxSanity, actual.MaxSanity)
	require.Equal(t, expected.CurrentSanity, actual.CurrentSanity)
	require.Equal(t, expected.TempInsanity, actual.TempInsanity)
	require.Equal(t, expected.IndefInsanity, actual.IndefInsanity)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualLuck(t *testing.T, expected db.LuckState, actual luckDTO.LuckModel) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CharacterID, actual.CharacterID)
	require.Equal(t, expected.StartingLuck, actual.StartingLuck)
	require.Equal(t, expected.CurrentLuck, actual.CurrentLuck)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func requireEqualNotes(t *testing.T, expected []db.Note, actual []notesDTO.NoteModel) {
	t.Helper()
	require.Len(t, actual, len(expected))
	for i, e := range expected {
		require.Equal(t, e.ID, actual[i].ID)
		require.Equal(t, e.CharacterID, actual[i].CharacterID)
		require.Equal(t, e.Title, actual[i].Title)
		require.Equal(t, e.Body, actual[i].Body)
		require.Equal(t, e.CreatedAt, actual[i].CreatedAt)
		require.Equal(t, e.UpdatedAt, actual[i].UpdatedAt)
	}
}
