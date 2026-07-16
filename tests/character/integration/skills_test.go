package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type defaultSkillDefinition struct {
	baseValue   int16
	isProtected bool
	baseRule    string
}

func defaultSkillDefinitions() map[string]defaultSkillDefinition {
	return map[string]defaultSkillDefinition{
		"Антропология":           {1, true, ""},
		"Археология":             {1, true, ""},
		"Ближний бой (драка)":    {25, true, ""},
		"Бухгалтерское дело":     {5, true, ""},
		"Верховая езда":          {5, true, ""},
		"Взлом":                  {1, true, ""},
		"Внимание":               {25, true, ""},
		"Вождение":               {20, true, ""},
		"Выживание":              {10, false, ""},
		"Естествознание":         {10, true, ""},
		"Запугивание":            {15, true, ""},
		"Искусство/ремесло":      {5, false, ""},
		"История":                {5, true, ""},
		"Красноречие":            {5, true, ""},
		"Лазание":                {20, true, ""},
		"Ловкость рук":           {10, true, ""},
		"Маскировка":             {5, true, ""},
		"Медицина":               {1, true, ""},
		"Метание":                {20, true, ""},
		"Механика":               {10, true, ""},
		"Мифы Ктулху":            {0, true, ""},
		"Наука":                  {1, false, ""},
		"Обаяние":                {15, true, ""},
		"Обоняние":               {15, true, ""},
		"Оккультизм":             {5, true, ""},
		"Ориентирование":         {10, true, ""},
		"Оценка":                 {5, true, ""},
		"Первая помощь":          {30, true, ""},
		"Пилотирование":          {1, true, ""},
		"Плавание":               {20, true, ""},
		"Прыжки":                 {20, true, ""},
		"Психоанализ":            {1, true, ""},
		"Психология":             {10, true, ""},
		"Работа в библиотеке":    {20, true, ""},
		"Скрытность":             {20, true, ""},
		"Слух":                   {20, true, ""},
		"Стрельба (винт./дроб.)": {25, true, ""},
		"Стрельба (пистолет)":    {20, true, ""},
		"Убеждение":              {10, true, ""},
		"Уклонение":              {0, true, "dodge"},
		"Управление тяжелыми машинами": {1, true, ""},
		"Чтение следов":                {10, true, ""},
		"Электрика":                    {10, true, ""},
		"Юриспруденция":                {5, true, ""},
		"Язык, иностранный":            {1, false, ""},
		"Язык, родной":                 {0, true, "native_language"},
	}
}

func requireDefaultSkills(t *testing.T, skills []db.Skill) {
	t.Helper()

	expected := defaultSkillDefinitions()
	require.Len(t, skills, len(expected))
	for _, skill := range skills {
		expectedSkill, ok := expected[skill.Name]
		require.True(t, ok, "unexpected skill %q", skill.Name)
		require.Equal(t, expectedSkill.baseValue, skill.BaseValue, skill.Name)
		require.Zero(t, skill.Value, skill.Name)
		require.False(t, skill.Checked, skill.Name)
		require.Equal(t, expectedSkill.isProtected, skill.IsProtected, skill.Name)
		if expectedSkill.baseRule == "" {
			require.Nil(t, skill.BaseRule, skill.Name)
		} else {
			require.NotNil(t, skill.BaseRule, skill.Name)
			require.Equal(t, expectedSkill.baseRule, *skill.BaseRule, skill.Name)
		}
	}
}

func TestCharacterServiceCreatesDefaultSkills(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: testUser.ID,
		Name:   "Default Skills Investigator",
	})
	require.NoError(t, err)

	skills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	requireDefaultSkills(t, skills)
}

func TestCharacterServiceSynchronizesDynamicSkillBasesAndPreservesImprovements(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: testUser.ID,
		Name:   "Dynamic Skills Investigator",
	})
	require.NoError(t, err)

	skills, err := service.GetSkills(context.Background(), skillsDTO.GetSkillsInput{UserID: testUser.ID, CharacterID: character.ID})
	require.NoError(t, err)
	var dodge skillsDTO.SkillModel
	for _, skill := range skills {
		if skill.Name == "Уклонение" {
			dodge = skill
			break
		}
	}
	require.True(t, dodge.ID.Valid)

	_, err = service.UpdateSkill(context.Background(), skillsDTO.UpdateSkillInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     dodge.ID,
		Name:        dodge.Name,
		BaseValue:   dodge.BaseValue,
		Value:       20,
		Checked:     true,
	})
	require.NoError(t, err)

	dexterity := int16(65)
	education := int16(80)
	_, err = service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Dexterity:   &dexterity,
		Education:   &education,
	})
	require.NoError(t, err)

	skills, err = service.GetSkills(context.Background(), skillsDTO.GetSkillsInput{UserID: testUser.ID, CharacterID: character.ID})
	require.NoError(t, err)
	byName := make(map[string]skillsDTO.SkillModel, len(skills))
	for _, skill := range skills {
		byName[skill.Name] = skill
	}
	require.Equal(t, int16(32), byName["Уклонение"].BaseValue)
	require.Equal(t, int16(20), byName["Уклонение"].Value)
	require.True(t, byName["Уклонение"].Checked)
	require.Equal(t, int16(80), byName["Язык, родной"].BaseValue)

	dexterity = 70
	education = 90
	_, err = service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Dexterity:   &dexterity,
		Education:   &education,
	})
	require.NoError(t, err)

	skills, err = service.GetSkills(context.Background(), skillsDTO.GetSkillsInput{UserID: testUser.ID, CharacterID: character.ID})
	require.NoError(t, err)
	byName = make(map[string]skillsDTO.SkillModel, len(skills))
	for _, skill := range skills {
		byName[skill.Name] = skill
	}
	require.Equal(t, int16(35), byName["Уклонение"].BaseValue)
	require.Equal(t, int16(20), byName["Уклонение"].Value)
	require.Equal(t, int16(90), byName["Язык, родной"].BaseValue)
}

func TestCharacterServiceProtectsFixedSkillsAndAllowsCustomSkillCRUD(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: testUser.ID,
		Name:   "Protected Skills Investigator",
	})
	require.NoError(t, err)

	skills, err := service.GetSkills(context.Background(), skillsDTO.GetSkillsInput{UserID: testUser.ID, CharacterID: character.ID})
	require.NoError(t, err)
	var fixed skillsDTO.SkillModel
	for _, skill := range skills {
		if skill.Name == "Антропология" {
			fixed = skill
			break
		}
	}
	require.True(t, fixed.ID.Valid)

	_, err = service.UpdateSkill(context.Background(), skillsDTO.UpdateSkillInput{
		UserID: testUser.ID, CharacterID: character.ID, SkillID: fixed.ID,
		Name: "Переименованная антропология", BaseValue: fixed.BaseValue,
	})
	require.ErrorIs(t, err, characterErrors.ErrProtectedSkill)

	_, err = service.UpdateSkill(context.Background(), skillsDTO.UpdateSkillInput{
		UserID: testUser.ID, CharacterID: character.ID, SkillID: fixed.ID,
		Name: fixed.Name, BaseValue: fixed.BaseValue + 1,
	})
	require.ErrorIs(t, err, characterErrors.ErrProtectedSkill)

	updatedFixed, err := service.UpdateSkill(context.Background(), skillsDTO.UpdateSkillInput{
		UserID: testUser.ID, CharacterID: character.ID, SkillID: fixed.ID,
		Name: fixed.Name, BaseValue: fixed.BaseValue, Value: 25, Checked: true,
	})
	require.NoError(t, err)
	require.Equal(t, int16(25), updatedFixed.Value)
	require.True(t, updatedFixed.Checked)

	err = service.DeleteSkill(context.Background(), skillsDTO.DeleteSkillInput{UserID: testUser.ID, CharacterID: character.ID, SkillID: fixed.ID})
	require.ErrorIs(t, err, characterErrors.ErrProtectedSkill)

	custom, err := service.CreateSkill(context.Background(), skillsDTO.CreateSkillInput{
		UserID: testUser.ID, CharacterID: character.ID,
		Name: "Наука: астрономия", BaseValue: 1, Value: 10,
	})
	require.NoError(t, err)
	require.False(t, custom.IsProtected)

	custom, err = service.UpdateSkill(context.Background(), skillsDTO.UpdateSkillInput{
		UserID: testUser.ID, CharacterID: character.ID, SkillID: custom.ID,
		Name: "Наука: астрофизика", BaseValue: 5, Value: 15, Checked: true,
	})
	require.NoError(t, err)
	require.Equal(t, "Наука: астрофизика", custom.Name)
	require.Equal(t, int16(5), custom.BaseValue)

	err = service.DeleteSkill(context.Background(), skillsDTO.DeleteSkillInput{UserID: testUser.ID, CharacterID: character.ID, SkillID: custom.ID})
	require.NoError(t, err)
}

func TestSkillsTableCreateListGetUpdateAndDeleteSkill(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdSkill, err := subject.queries.CreateCharacterSkill(context.Background(), db.CreateCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Name:        "Library Use",
		BaseValue:   20,
		Value:       45,
		Checked:     true,
	})
	require.NoError(t, err)

	require.True(t, createdSkill.ID.Valid)
	require.Equal(t, character.ID, createdSkill.CharacterID)
	require.Equal(t, "Library Use", createdSkill.Name)
	require.Equal(t, int16(20), createdSkill.BaseValue)
	require.Equal(t, int16(45), createdSkill.Value)
	require.True(t, createdSkill.Checked)
	require.False(t, createdSkill.IsProtected)
	require.Nil(t, createdSkill.BaseRule)

	fetchedSkill, err := subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     createdSkill.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdSkill.ID, fetchedSkill.ID)

	time.Sleep(5 * time.Millisecond)

	updatedSkill, err := subject.queries.UpdateCharacterSkill(context.Background(), db.UpdateCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     createdSkill.ID,
		Name:        "Library Use Updated",
		BaseValue:   25,
		Value:       60,
		Checked:     false,
	})
	require.NoError(t, err)

	require.Equal(t, createdSkill.ID, updatedSkill.ID)
	require.Equal(t, "Library Use Updated", updatedSkill.Name)
	require.Equal(t, int16(25), updatedSkill.BaseValue)
	require.Equal(t, int16(60), updatedSkill.Value)
	require.False(t, updatedSkill.Checked)
	require.False(t, updatedSkill.IsProtected)
	require.Nil(t, updatedSkill.BaseRule)
	require.True(t, updatedSkill.UpdatedAt.Time.After(createdSkill.UpdatedAt.Time) || updatedSkill.UpdatedAt.Time.Equal(createdSkill.UpdatedAt.Time))

	deletedSkill, err := subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     createdSkill.ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdSkill.ID, deletedSkill.ID)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     createdSkill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestSkillsTableListReturnsSkillsOrderedByName(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Zoology"))
	require.NoError(t, err)
	_, err = subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Accounting"))
	require.NoError(t, err)

	skills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, skills, 2)
	require.Equal(t, "Accounting", skills[0].Name)
	require.Equal(t, "Zoology", skills[1].Name)
}

func TestSkillsTableListReturnsEmptyForCharacterWithoutSkills(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, skills)
}

func TestSkillsTableAllowsZeroValues(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	params := testCreateSkillParams(testUser.ID, character.ID, "Zero Skill")
	params.BaseValue = 0
	params.Value = 0

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, int16(0), skill.BaseValue)
	require.Equal(t, int16(0), skill.Value)
}

func TestSkillsTableAllowsEmptyName(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, ""))
	require.NoError(t, err)
	require.Equal(t, "", skill.Name)

	updatedSkill, err := subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, skill.ID, ""))
	require.NoError(t, err)
	require.Equal(t, "", updatedSkill.Name)
}

func TestSkillsTableRequiresCharacterOwnerForCreateListGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(otherUser.ID, character.ID, "Unauthorized Skill"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(owner.ID, character.ID, "Owner Skill"))
	require.NoError(t, err)

	skills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, skills)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(otherUser.ID, character.ID, skill.ID, "Unauthorized Update"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestSkillsTableUnauthorizedUpdateDoesNotMutateExistingSkill(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(owner.ID, character.ID, "Original Skill"))
	require.NoError(t, err)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(otherUser.ID, character.ID, skill.ID, "Unauthorized Mutation"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedSkill, err := subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.NoError(t, err)
	require.Equal(t, "Original Skill", fetchedSkill.Name)
	require.Equal(t, int16(10), fetchedSkill.BaseValue)
	require.Equal(t, int16(35), fetchedSkill.Value)
}

func TestSkillsTableRequiresMatchingCharacterForGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	owningCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	otherCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, owningCharacter.ID, "Scoped Skill"))
	require.NoError(t, err)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, otherCharacter.ID, skill.ID, "Wrong Character"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestSkillsTableKeepsSkillsScopedToRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, firstCharacter.ID, "First Skill"))
	require.NoError(t, err)
	secondSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, secondCharacter.ID, "Second Skill"))
	require.NoError(t, err)

	firstSkills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Len(t, firstSkills, 1)
	require.Equal(t, firstSkill.ID, firstSkills[0].ID)

	secondSkills, err := subject.queries.GetCharacterSkills(context.Background(), db.GetCharacterSkillsParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Len(t, secondSkills, 1)
	require.Equal(t, secondSkill.ID, secondSkills[0].ID)
}

func TestSkillsTableReturnsNoRowsForMissingCharacterOrSkill(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	missingCharacterID := characterTestUUID("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	missingSkillID := characterTestUUID("ffffffff-ffff-ffff-ffff-ffffffffffff")

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, missingCharacterID, "Missing Character"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     missingSkillID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, missingSkillID, "Missing Skill"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     missingSkillID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestSkillsTableRejectsNegativeValues(t *testing.T) {
	tests := []struct {
		name   string
		params func(userID string, characterID pgtype.UUID) db.CreateCharacterSkillParams
	}{
		{
			name: "base value",
			params: func(userID string, characterID pgtype.UUID) db.CreateCharacterSkillParams {
				params := testCreateSkillParams(userID, characterID, "Negative Base")
				params.BaseValue = -1
				return params
			},
		},
		{
			name: "value",
			params: func(userID string, characterID pgtype.UUID) db.CreateCharacterSkillParams {
				params := testCreateSkillParams(userID, characterID, "Negative Value")
				params.Value = -1
				return params
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)

			_, err := subject.queries.CreateCharacterSkill(context.Background(), tc.params(testUser.ID, character.ID))
			requirePostgresErrorCode(t, err, "23514")
		})
	}
}

func TestSkillsTableRejectsInvalidUpdateValues(t *testing.T) {
	tests := []struct {
		name   string
		params func(userID string, characterID pgtype.UUID, skillID pgtype.UUID) db.UpdateCharacterSkillParams
	}{
		{
			name: "negative base value",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, "Negative Base Update")
				params.BaseValue = -1
				return params
			},
		},
		{
			name: "negative value",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, "Negative Value Update")
				params.Value = -1
				return params
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)
			skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Valid Skill"))
			require.NoError(t, err)

			_, err = subject.queries.UpdateCharacterSkill(context.Background(), tc.params(testUser.ID, character.ID, skill.ID))
			requirePostgresErrorCode(t, err, "23514")
		})
	}
}

func TestSkillsTableRejectsNameLongerThanLimit(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	longName := strings.Repeat("a", 101)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, longName))
	requirePostgresErrorCode(t, err, "22001")

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Valid Skill"))
	require.NoError(t, err)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, skill.ID, longName))
	requirePostgresErrorCode(t, err, "22001")
}

func TestSkillsTableDeleteReturnsDeletedValuesAndAllowsRecreate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Delete Me"))
	require.NoError(t, err)

	deletedSkill, err := subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.NoError(t, err)
	require.Equal(t, skill.ID, deletedSkill.ID)
	require.Equal(t, "Delete Me", deletedSkill.Name)

	recreatedSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Delete Me"))
	require.NoError(t, err)
	require.NotEqual(t, deletedSkill.ID, recreatedSkill.ID)
}

func TestSkillsTableDeletingCharacterCascadesSkills(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, "Cascade Skill"))
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
