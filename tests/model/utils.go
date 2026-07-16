package tests

import (
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func testUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	if err != nil {
		panic(err)
	}

	return uuid
}

func invalidUUID() pgtype.UUID {
	return pgtype.UUID{}
}

func testTimestamptz(value string) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	err := ts.Scan(value)
	if err != nil {
		panic(err)
	}

	return ts
}

func strPtr(value string) *string {
	return &value
}

func int16Ptr(value int16) *int16 {
	return &value
}

func testCharacter() db.Character {
	portraitKey := "portraits/11111111-1111-1111-1111-111111111111.webp"
	return db.Character{
		ID:          testUUID("11111111-1111-1111-1111-111111111111"),
		UserID:      "user_1",
		Name:        "Dr. Armitage",
		Occupation:  strPtr("Antiquarian"),
		Age:         int16Ptr(42),
		Sex:         strPtr("male"),
		Residence:   strPtr("Arkham"),
		Birthplace:  strPtr("Boston"),
		CreatedAt:   testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:   testTimestamptz("2026-06-07 13:00:00+03"),
		PortraitKey: &portraitKey,
	}
}

func testBackstory() db.Backstory {
	return db.Backstory{
		ID:                  testUUID("22222222-2222-2222-2222-222222222222"),
		CharacterID:         testUUID("11111111-1111-1111-1111-111111111111"),
		PersonalDescription: strPtr("A careful scholar with tired eyes."),
		CreatedAt:           testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:           testTimestamptz("2026-06-07 13:00:00+03"),
	}
}

func testBackstoryItem() db.BackstoryItem {
	return db.BackstoryItem{
		ID:          testUUID("33333333-3333-3333-3333-333333333333"),
		BackstoryID: testUUID("22222222-2222-2222-2222-222222222222"),
		Section:     "ideology_beliefs",
		Title:       "The old motto",
		Text:        "Knowledge has a price.",
		CreatedAt:   testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:   testTimestamptz("2026-06-07 13:00:00+03"),
	}
}

func testSkillRow() db.Skill {
	return db.Skill{
		ID:          testUUID("44444444-4444-4444-4444-444444444444"),
		CharacterID: testUUID("11111111-1111-1111-1111-111111111111"),
		Name:        "Library Use",
		BaseValue:   20,
		Value:       65,
		Checked:     true,
		IsProtected: true,
		BaseRule:    strPtr("dodge"),
		CreatedAt:   testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:   testTimestamptz("2026-06-07 13:00:00+03"),
	}
}

func testFinance() db.Finance {
	return db.Finance{
		ID:                  testUUID("88888888-8888-8888-8888-888888888888"),
		CharacterID:         testUUID("11111111-1111-1111-1111-111111111111"),
		SpendingLimit:       strPtr("$50"),
		Cash:                strPtr("$120"),
		Assets:              strPtr("Books and a battered Ford."),
		CreditRatingSkillID: invalidUUID(),
		CreatedAt:           testTimestamptz("2026-06-07 12:00:00+03"),
		UpdatedAt:           testTimestamptz("2026-06-07 13:00:00+03"),
	}
}

func requireSameShortCharacter(t *testing.T, expected db.Character, actual characterDTO.CharacterShortModel) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.UserID, actual.UserID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Occupation, actual.Occupation)
	require.Equal(t, expected.Age, actual.Age)
	require.Equal(t, expected.Sex, actual.Sex)
	require.Equal(t, expected.Residence, actual.Residence)
	require.Equal(t, expected.Birthplace, actual.Birthplace)
	require.Nil(t, actual.PortraitUrl)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}

func requireSameSkill(t *testing.T, expected db.Skill, actual skillsDTO.SkillModel) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.BaseValue, actual.BaseValue)
	require.Equal(t, expected.Value, actual.Value)
	require.Equal(t, int32(expected.BaseValue)+int32(expected.Value), actual.TotalValue)
	require.Equal(t, expected.Checked, actual.Checked)
	require.Equal(t, expected.IsProtected, actual.IsProtected)
	require.Equal(t, expected.BaseRule, actual.BaseRule)
	require.Equal(t, expected.CreatedAt.Time, actual.CreatedAt.Time)
	require.Equal(t, expected.UpdatedAt.Time, actual.UpdatedAt.Time)
}
