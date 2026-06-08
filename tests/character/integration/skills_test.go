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

func TestSkillsTableCreateListGetUpdateAndDeleteSkill(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, categoryName := createSkillTestCategory(t, subject, "Investigation")
	specialtyID, specialtyName := createSkillTestSpecialty(t, subject, "Library Use", "Research in archives.", 20)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	createdSkill, err := subject.queries.CreateCharacterSkill(context.Background(), db.CreateCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Name:        "Library Use",
		CategoryID:  categoryID,
		BaseValue:   20,
		Value:       45,
		Checked:     true,
		Specialized: true,
		SpecialtyID: specialtyID,
	})
	require.NoError(t, err)

	require.True(t, createdSkill.ID.Valid)
	require.Equal(t, character.ID, createdSkill.CharacterID)
	require.Equal(t, "Library Use", createdSkill.Name)
	require.Equal(t, categoryID, createdSkill.CategoryID)
	require.Equal(t, int16(20), createdSkill.BaseValue)
	require.Equal(t, int16(45), createdSkill.Value)
	require.True(t, createdSkill.Checked)
	require.True(t, createdSkill.Specialized)
	require.Equal(t, specialtyID, createdSkill.SpecialtyID)
	require.Equal(t, categoryName, createdSkill.CategoryName)
	require.Equal(t, specialtyID, createdSkill.SpecialtyPkID)
	require.NotNil(t, createdSkill.SpecialtyName)
	require.Equal(t, specialtyName, *createdSkill.SpecialtyName)
	require.NotNil(t, createdSkill.SpecialtyDescription)
	require.Equal(t, "Research in archives.", *createdSkill.SpecialtyDescription)
	require.NotNil(t, createdSkill.SpecialtyBaseValue)
	require.Equal(t, int16(20), *createdSkill.SpecialtyBaseValue)

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
		CategoryID:  categoryID,
		BaseValue:   25,
		Value:       60,
		Checked:     false,
		Specialized: false,
		SpecialtyID: pgtype.UUID{},
	})
	require.NoError(t, err)

	require.Equal(t, createdSkill.ID, updatedSkill.ID)
	require.Equal(t, "Library Use Updated", updatedSkill.Name)
	require.Equal(t, int16(25), updatedSkill.BaseValue)
	require.Equal(t, int16(60), updatedSkill.Value)
	require.False(t, updatedSkill.Checked)
	require.False(t, updatedSkill.Specialized)
	require.False(t, updatedSkill.SpecialtyID.Valid)
	require.Nil(t, updatedSkill.SpecialtyName)
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
	categoryID, _ := createSkillTestCategory(t, subject, "General")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Zoology"))
	require.NoError(t, err)
	_, err = subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Accounting"))
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
	categoryID, _ := createSkillTestCategory(t, subject, "Zero Values")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	params := testCreateSkillParams(testUser.ID, character.ID, categoryID, "Zero Skill")
	params.BaseValue = 0
	params.Value = 0

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, int16(0), skill.BaseValue)
	require.Equal(t, int16(0), skill.Value)
}

func TestSkillsTableAllowsEmptyName(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Empty Name")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, ""))
	require.NoError(t, err)
	require.Equal(t, "", skill.Name)

	updatedSkill, err := subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, skill.ID, categoryID, ""))
	require.NoError(t, err)
	require.Equal(t, "", updatedSkill.Name)
}

func TestSkillsTableRequiresCharacterOwnerForCreateListGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Owner Scoped")
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(otherUser.ID, character.ID, categoryID, "Unauthorized Skill"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(owner.ID, character.ID, categoryID, "Owner Skill"))
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

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(otherUser.ID, character.ID, skill.ID, categoryID, "Unauthorized Update"))
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
	categoryID, _ := createSkillTestCategory(t, subject, "Unauthorized Mutation")
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(owner.ID, character.ID, categoryID, "Original Skill"))
	require.NoError(t, err)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(otherUser.ID, character.ID, skill.ID, categoryID, "Unauthorized Mutation"))
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
	categoryID, _ := createSkillTestCategory(t, subject, "Character Scoped")
	testUser := createCharacterTestUser(t, subject)
	owningCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	otherCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, owningCharacter.ID, categoryID, "Scoped Skill"))
	require.NoError(t, err)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		SkillID:     skill.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, otherCharacter.ID, skill.ID, categoryID, "Wrong Character"))
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
	categoryID, _ := createSkillTestCategory(t, subject, "Multi Character")
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, firstCharacter.ID, categoryID, "First Skill"))
	require.NoError(t, err)
	secondSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, secondCharacter.ID, categoryID, "Second Skill"))
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
	categoryID, _ := createSkillTestCategory(t, subject, "Missing Rows")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	missingCharacterID := characterTestUUID("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	missingSkillID := characterTestUUID("ffffffff-ffff-ffff-ffff-ffffffffffff")

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, missingCharacterID, categoryID, "Missing Character"))
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetCharacterSkill(context.Background(), db.GetCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     missingSkillID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, missingSkillID, categoryID, "Missing Skill"))
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
		params func(userID string, characterID pgtype.UUID, categoryID pgtype.UUID) db.CreateCharacterSkillParams
	}{
		{
			name: "base value",
			params: func(userID string, characterID pgtype.UUID, categoryID pgtype.UUID) db.CreateCharacterSkillParams {
				params := testCreateSkillParams(userID, characterID, categoryID, "Negative Base")
				params.BaseValue = -1
				return params
			},
		},
		{
			name: "value",
			params: func(userID string, characterID pgtype.UUID, categoryID pgtype.UUID) db.CreateCharacterSkillParams {
				params := testCreateSkillParams(userID, characterID, categoryID, "Negative Value")
				params.Value = -1
				return params
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			categoryID, _ := createSkillTestCategory(t, subject, "Negative "+tc.name)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)

			_, err := subject.queries.CreateCharacterSkill(context.Background(), tc.params(testUser.ID, character.ID, categoryID))
			requirePostgresErrorCode(t, err, "23514")
		})
	}
}

func TestSkillsTableRejectsMissingCategoryOrSpecialty(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "FK Category")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	missingCategoryID := characterTestUUID("12121212-1212-1212-1212-121212121212")
	missingSpecialtyID := characterTestUUID("34343434-3434-3434-3434-343434343434")

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, missingCategoryID, "Missing Category"))
	requirePostgresErrorCode(t, err, "23503")

	params := testCreateSkillParams(testUser.ID, character.ID, categoryID, "Missing Specialty")
	params.Specialized = true
	params.SpecialtyID = missingSpecialtyID
	_, err = subject.queries.CreateCharacterSkill(context.Background(), params)
	requirePostgresErrorCode(t, err, "23503")
}

func TestSkillsTableRejectsInvalidUpdateValues(t *testing.T) {
	tests := []struct {
		name   string
		params func(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID) db.UpdateCharacterSkillParams
	}{
		{
			name: "negative base value",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, categoryID, "Negative Base Update")
				params.BaseValue = -1
				return params
			},
		},
		{
			name: "negative value",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, categoryID, "Negative Value Update")
				params.Value = -1
				return params
			},
		},
		{
			name: "missing category",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, categoryID, "Missing Category Update")
				params.CategoryID = characterTestUUID("56565656-5656-5656-5656-565656565656")
				return params
			},
		},
		{
			name: "missing specialty",
			params: func(userID string, characterID pgtype.UUID, skillID pgtype.UUID, categoryID pgtype.UUID) db.UpdateCharacterSkillParams {
				params := testUpdateSkillParams(userID, characterID, skillID, categoryID, "Missing Specialty Update")
				params.Specialized = true
				params.SpecialtyID = characterTestUUID("78787878-7878-7878-7878-787878787878")
				return params
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subject := newCharacterIntegrationSubject(t)
			categoryID, _ := createSkillTestCategory(t, subject, "Invalid Update "+tc.name)
			testUser := createCharacterTestUser(t, subject)
			character := createCharacterTestCharacter(t, subject, testUser.ID)
			skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Valid Skill"))
			require.NoError(t, err)

			_, err = subject.queries.UpdateCharacterSkill(context.Background(), tc.params(testUser.ID, character.ID, skill.ID, categoryID))
			if strings.Contains(tc.name, "negative") {
				requirePostgresErrorCode(t, err, "23514")
				return
			}
			requirePostgresErrorCode(t, err, "23503")
		})
	}
}

func TestSkillsTableRejectsNameLongerThanLimit(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Name Limit")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	longName := strings.Repeat("a", 101)

	_, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, longName))
	requirePostgresErrorCode(t, err, "22001")

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Valid Skill"))
	require.NoError(t, err)

	_, err = subject.queries.UpdateCharacterSkill(context.Background(), testUpdateSkillParams(testUser.ID, character.ID, skill.ID, categoryID, longName))
	requirePostgresErrorCode(t, err, "22001")
}

func TestSkillsTableAllowsSpecializedWithoutSpecialtyAndSpecialtyWhenNotSpecialized(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Specialization Flags")
	specialtyID, _ := createSkillTestSpecialty(t, subject, "Flag Specialty", "Specialty flag edge cases.", 5)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	specializedWithoutSpecialtyParams := testCreateSkillParams(testUser.ID, character.ID, categoryID, "Specialized Without Specialty")
	specializedWithoutSpecialtyParams.Specialized = true
	specializedWithoutSpecialtyParams.SpecialtyID = pgtype.UUID{}
	specializedWithoutSpecialty, err := subject.queries.CreateCharacterSkill(context.Background(), specializedWithoutSpecialtyParams)
	require.NoError(t, err)
	require.True(t, specializedWithoutSpecialty.Specialized)
	require.False(t, specializedWithoutSpecialty.SpecialtyID.Valid)
	require.Nil(t, specializedWithoutSpecialty.SpecialtyName)

	specialtyWhenNotSpecializedParams := testCreateSkillParams(testUser.ID, character.ID, categoryID, "Specialty When Not Specialized")
	specialtyWhenNotSpecializedParams.Specialized = false
	specialtyWhenNotSpecializedParams.SpecialtyID = specialtyID
	specialtyWhenNotSpecialized, err := subject.queries.CreateCharacterSkill(context.Background(), specialtyWhenNotSpecializedParams)
	require.NoError(t, err)
	require.False(t, specialtyWhenNotSpecialized.Specialized)
	require.Equal(t, specialtyID, specialtyWhenNotSpecialized.SpecialtyID)
	require.NotNil(t, specialtyWhenNotSpecialized.SpecialtyName)
}

func TestSkillsTableDeleteReturnsDeletedValuesAndAllowsRecreate(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Delete Recreate")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Delete Me"))
	require.NoError(t, err)

	deletedSkill, err := subject.queries.DeleteCharacterSkill(context.Background(), db.DeleteCharacterSkillParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		SkillID:     skill.ID,
	})
	require.NoError(t, err)
	require.Equal(t, skill.ID, deletedSkill.ID)
	require.Equal(t, "Delete Me", deletedSkill.Name)

	recreatedSkill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Delete Me"))
	require.NoError(t, err)
	require.NotEqual(t, deletedSkill.ID, recreatedSkill.ID)
}

func TestSkillsTableDeletingCharacterCascadesSkills(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	categoryID, _ := createSkillTestCategory(t, subject, "Cascade")
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	skill, err := subject.queries.CreateCharacterSkill(context.Background(), testCreateSkillParams(testUser.ID, character.ID, categoryID, "Cascade Skill"))
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
