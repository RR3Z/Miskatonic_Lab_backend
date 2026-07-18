package tests

import (
	"context"
	"strings"
	"testing"

	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateSkillRejectsBlankName(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateSkill(context.Background(), skillsDTO.CreateSkillInput{Name: "   ", BaseValue: 1, Value: 1})
	require.ErrorIs(t, err, characterErrors.ErrSkillNameRequired)
}

func TestCreateSkillRejectsNameTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longName := string(make([]byte, 101))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	_, err := service.CreateSkill(context.Background(), skillsDTO.CreateSkillInput{Name: longName, BaseValue: 1, Value: 1})
	require.ErrorIs(t, err, characterErrors.ErrSkillNameTooLong)
}

func TestCreateSkillRejectsNegativeBaseValue(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateSkill(context.Background(), skillsDTO.CreateSkillInput{Name: "test", BaseValue: -1, Value: 1})
	require.ErrorIs(t, err, characterErrors.ErrSkillValueNegative)
}

func TestCreateSkillRejectsNegativeValue(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateSkill(context.Background(), skillsDTO.CreateSkillInput{Name: "test", BaseValue: 1, Value: -1})
	require.ErrorIs(t, err, characterErrors.ErrSkillValueNegative)
}

func TestCreateBackstoryItemRejectsBlankSection(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateBackstoryItem(context.Background(), backstoriesDTO.CreateBackstoryItemInput{Section: "   ", Title: "title", Text: "text"})
	require.ErrorIs(t, err, characterErrors.ErrInvalidBackstorySection)
}

func TestCreateBackstoryItemRejectsSectionTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longSection := string(make([]byte, 33))
	for i := range longSection {
		longSection = longSection[:i] + "a" + longSection[i+1:]
	}
	_, err := service.CreateBackstoryItem(context.Background(), backstoriesDTO.CreateBackstoryItemInput{Section: longSection, Title: "title", Text: "text"})
	require.ErrorIs(t, err, characterErrors.ErrSectionTooLong)
}

func TestCreateBackstoryItemRejectsBlankTitle(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateBackstoryItem(context.Background(), backstoriesDTO.CreateBackstoryItemInput{Section: "injuries_scars", Title: "   ", Text: "text"})
	require.ErrorIs(t, err, characterErrors.ErrBackstoryTitleRequired)
}

func TestCreateBackstoryItemRejectsTitleTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longTitle := string(make([]byte, 256))
	for i := range longTitle {
		longTitle = longTitle[:i] + "a" + longTitle[i+1:]
	}
	_, err := service.CreateBackstoryItem(context.Background(), backstoriesDTO.CreateBackstoryItemInput{Section: "injuries_scars", Title: longTitle, Text: "text"})
	require.ErrorIs(t, err, characterErrors.ErrBackstoryTitleTooLong)
}

func TestCreateBackstoryItemRejectsBlankText(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateBackstoryItem(context.Background(), backstoriesDTO.CreateBackstoryItemInput{Section: "injuries_scars", Title: "title", Text: "   "})
	require.ErrorIs(t, err, characterErrors.ErrBackstoryTextRequired)
}

func TestCreateNoteRejectsBlankTitle(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateNote(context.Background(), notesDTO.CreateNoteInput{Title: "   ", Body: "body"})
	require.ErrorIs(t, err, characterErrors.ErrNoteTitleRequired)
}

func TestCreateNoteRejectsTitleTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longTitle := string(make([]byte, 121))
	for i := range longTitle {
		longTitle = longTitle[:i] + "a" + longTitle[i+1:]
	}
	_, err := service.CreateNote(context.Background(), notesDTO.CreateNoteInput{Title: longTitle, Body: "body"})
	require.ErrorIs(t, err, characterErrors.ErrNoteTitleTooLong)
}

func TestCreateNoteRejectsBlankBody(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateNote(context.Background(), notesDTO.CreateNoteInput{Title: "title", Body: "   "})
	require.ErrorIs(t, err, characterErrors.ErrNoteBodyRequired)
}

func TestCreateInventoryItemRejectsBlankName(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateInventoryItem(context.Background(), inventoryDTO.CreateInventoryItemInput{Name: "   "})
	require.ErrorIs(t, err, characterErrors.ErrInventoryItemNameRequired)
}

func TestCreateInventoryItemRejectsTooLongName(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	_, err := service.CreateInventoryItem(context.Background(), inventoryDTO.CreateInventoryItemInput{Name: strings.Repeat("a", 121)})
	require.ErrorIs(t, err, characterErrors.ErrInventoryItemNameTooLong)
}

func TestCreateInventoryItemRejectsNonPositiveQuantity(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	quantity := int32(0)
	_, err := service.CreateInventoryItem(context.Background(), inventoryDTO.CreateInventoryItemInput{Name: "Flashlight", Quantity: &quantity})
	require.ErrorIs(t, err, characterErrors.ErrInventoryItemQuantityInvalid)
}

func TestCreateInventoryItemRejectsTooLongCategory(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	category := strings.Repeat("a", 81)
	_, err := service.CreateInventoryItem(context.Background(), inventoryDTO.CreateInventoryItemInput{Name: "Flashlight", Category: &category})
	require.ErrorIs(t, err, characterErrors.ErrInventoryItemCategoryTooLong)
}

func TestUpsertCharacteristicsRejectsNegativeValue(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	negVal := int16(-1)
	_, err := service.UpsertCharacteristics(context.Background(), characteristicsDTO.UpsertCharacteristicsInput{
		Strength: &negVal,
	})
	require.ErrorIs(t, err, characterErrors.ErrCharacteristicsNegative)
}

func TestUpsertFinancesRejectsSpendingLimitTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longVal := string(make([]byte, 121))
	for i := range longVal {
		longVal = longVal[:i] + "a" + longVal[i+1:]
	}
	_, err := service.UpsertFinances(context.Background(), financesDTO.UpsertFinancesInput{
		SpendingLimit: &longVal,
	})
	require.ErrorIs(t, err, characterErrors.ErrFinancesMoneyTooLong)
}

func TestUpsertFinancesRejectsCashTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longVal := string(make([]byte, 121))
	for i := range longVal {
		longVal = longVal[:i] + "a" + longVal[i+1:]
	}
	_, err := service.UpsertFinances(context.Background(), financesDTO.UpsertFinancesInput{
		Cash: &longVal,
	})
	require.ErrorIs(t, err, characterErrors.ErrFinancesMoneyTooLong)
}
