package character

import (
	"context"
	"errors"
	"strings"

	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

func validateRequiredString(s string, maxLen int, requiredErr, tooLongErr error) error {
	if strings.TrimSpace(s) == "" {
		return requiredErr
	}
	if maxLen > 0 && len(s) > maxLen {
		return tooLongErr
	}
	return nil
}

func validateNonNegative(err error, fields ...*int16) error {
	for _, f := range fields {
		if f != nil && *f < 0 {
			return err
		}
	}
	return nil
}

func validateSex(sex *string) error {
	if sex == nil {
		return nil
	}

	switch *sex {
	case "male", "female":
		return nil
	default:
		return characterErrors.ErrSexInvalid
	}
}

func validateSkillInput(name string, baseValue, value int16) error {
	if err := validateRequiredString(name, 100, characterErrors.ErrSkillNameRequired, characterErrors.ErrSkillNameTooLong); err != nil {
		return err
	}
	if baseValue < 0 || value < 0 {
		return characterErrors.ErrSkillValueNegative
	}
	return nil
}

func validateBackstoryItemInput(section, title, text string) error {
	if err := validateRequiredString(section, 32, characterErrors.ErrInvalidBackstorySection, characterErrors.ErrSectionTooLong); err != nil {
		return err
	}
	if err := validateRequiredString(title, 255, characterErrors.ErrBackstoryTitleRequired, characterErrors.ErrBackstoryTitleTooLong); err != nil {
		return err
	}
	if err := validateRequiredString(text, 0, characterErrors.ErrBackstoryTextRequired, nil); err != nil {
		return err
	}
	return nil
}

func validateNoteInput(title, body string) error {
	if err := validateRequiredString(title, 120, characterErrors.ErrNoteTitleRequired, characterErrors.ErrNoteTitleTooLong); err != nil {
		return err
	}
	if err := validateRequiredString(body, 0, characterErrors.ErrNoteBodyRequired, nil); err != nil {
		return err
	}
	return nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func normalizeInventoryItemInput(name string, quantity *int32, category, description *string) (string, *int32, *string, *string, error) {
	normalizedName := strings.TrimSpace(name)
	if err := validateRequiredString(normalizedName, 120, characterErrors.ErrInventoryItemNameRequired, characterErrors.ErrInventoryItemNameTooLong); err != nil {
		return "", nil, nil, nil, err
	}
	if quantity != nil && *quantity < 1 {
		return "", nil, nil, nil, characterErrors.ErrInventoryItemQuantityInvalid
	}

	normalizedCategory := normalizeOptionalString(category)
	if normalizedCategory != nil && len(*normalizedCategory) > 80 {
		return "", nil, nil, nil, characterErrors.ErrInventoryItemCategoryTooLong
	}

	return normalizedName, quantity, normalizedCategory, normalizeOptionalString(description), nil
}

type stateFieldFetcher func(ctx context.Context) (max int16, current int16, err error)

func (s *CharacterService) validateStateMax(
	ctx context.Context,
	maxVal, currentVal *int16,
	fetchExisting stateFieldFetcher,
	exceedsErr error,
) error {
	if maxVal != nil && currentVal != nil {
		if *currentVal > *maxVal {
			return exceedsErr
		}
		return nil
	}

	if maxVal == nil && currentVal == nil {
		return nil
	}

	existingMax, existingCurrent, err := fetchExisting(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	if maxVal != nil {
		existingMax = *maxVal
	}
	if currentVal != nil {
		existingCurrent = *currentVal
	}

	if existingCurrent > existingMax {
		return exceedsErr
	}

	return nil
}
