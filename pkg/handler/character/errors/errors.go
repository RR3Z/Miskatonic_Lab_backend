package characterErrors

import (
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, characterErrors.ErrNameRequired):
		return characterNameRequiredError(err)
	case errors.Is(err, characterErrors.ErrNameTooLong):
		return badRequestError("character.name_too_long", "name exceeds maximum length", err)
	case errors.Is(err, characterErrors.ErrAgeNegative):
		return badRequestError("character.age_negative", "age must be >= 0", err)
	case errors.Is(err, characterErrors.ErrSexInvalid):
		return badRequestError("character.sex_invalid", "sex must be male or female", err)
	case errors.Is(err, characterErrors.ErrPatchInvalid):
		return InvalidInputError("character patch input is invalid", err)
	case errors.Is(err, characterErrors.ErrCharacterLimitReached):
		return &myErrors.AppError{
			Status:  http.StatusConflict,
			Code:    "character.limit_reached",
			Message: "maximum number of characters reached",
			Details: []myErrors.ErrorDetail{
				myErrors.ConflictDetail("characters", "limit_reached"),
			},
			Err: err,
		}
	case errors.Is(err, characterErrors.ErrPortraitRequired):
		return badRequestError("character.portrait_required", "portrait file is required", err)
	case errors.Is(err, characterErrors.ErrPortraitTooLarge):
		return &myErrors.AppError{Status: http.StatusRequestEntityTooLarge, Code: "character.portrait_too_large", Message: "portrait exceeds 5 MiB", Err: err}
	case errors.Is(err, characterErrors.ErrPortraitUnsupported):
		return badRequestError("character.portrait_unsupported", "portrait must be JPEG, PNG, or WebP", err)
	case errors.Is(err, characterErrors.ErrPortraitInvalid):
		return badRequestError("character.portrait_invalid", "portrait must be a valid image up to 4096x4096", err)
	case errors.Is(err, characterErrors.ErrPortraitStorage):
		return &myErrors.AppError{Status: http.StatusServiceUnavailable, Code: "character.portrait_storage_unavailable", Message: "portrait storage is unavailable", Err: err}
	case errors.Is(err, characterErrors.ErrCharacteristicsNegative):
		return badRequestError("character.characteristics_negative", "characteristic values must be >= 0", err)
	case errors.Is(err, characterErrors.ErrDerivedStatsNegative):
		return badRequestError("character.derived_stats_negative", "speed and dodge must be >= 0", err)
	case errors.Is(err, characterErrors.ErrInvalidDamageBonus):
		return badRequestError("character.invalid_damage_bonus", "invalid damage bonus format", err)
	case errors.Is(err, characterErrors.ErrStateNegative):
		return badRequestError("character.state_negative", "state values must be >= 0", err)
	case errors.Is(err, characterErrors.ErrSectionTooLong):
		return badRequestError("character.section_too_long", "backstory item section exceeds max length", err)
	case errors.Is(err, characterErrors.ErrBackstoryTitleRequired):
		return badRequestError("character.backstory_title_required", "backstory item title is required", err)
	case errors.Is(err, characterErrors.ErrBackstoryTitleTooLong):
		return badRequestError("character.backstory_title_too_long", "backstory item title exceeds max length", err)
	case errors.Is(err, characterErrors.ErrBackstoryTextRequired):
		return badRequestError("character.backstory_text_required", "backstory item text is required", err)
	case errors.Is(err, characterErrors.ErrSkillNameRequired):
		return badRequestError("character.skill_name_required", "skill name is required", err)
	case errors.Is(err, characterErrors.ErrSkillNameTooLong):
		return badRequestError("character.skill_name_too_long", "skill name exceeds max length", err)
	case errors.Is(err, characterErrors.ErrSkillValueNegative):
		return badRequestError("character.skill_value_negative", "skill values must be >= 0", err)
	case errors.Is(err, characterErrors.ErrProtectedSkill):
		return &myErrors.AppError{Status: http.StatusConflict, Code: "character.skill_protected", Message: "protected skill cannot be renamed, rebased, or deleted", Err: err}
	case errors.Is(err, characterErrors.ErrFinancesMoneyTooLong):
		return badRequestError("character.finances_money_too_long", "money field exceeds max length", err)
	case errors.Is(err, characterErrors.ErrNoteTitleRequired):
		return badRequestError("character.note_title_required", "note title is required", err)
	case errors.Is(err, characterErrors.ErrNoteTitleTooLong):
		return badRequestError("character.note_title_too_long", "note title exceeds max length", err)
	case errors.Is(err, characterErrors.ErrNoteBodyRequired):
		return badRequestError("character.note_body_required", "note body is required", err)
	case isHealthStateValidationError(err):
		return badRequestError("character.state_current_exceeds_max", "current_hp value cannot exceed max_hp value", err)
	case isMagicStateValidationError(err):
		return badRequestError("character.state_current_exceeds_max", "current_mp value cannot exceed max_mp value", err)
	case isSanityStateValidationError(err):
		return badRequestError("character.state_current_exceeds_max", "current_sanity value cannot exceed max_sanity value", err)
	case isLuckStateValidationError(err):
		return badRequestError("character.state_current_exceeds_max", "current_luck value cannot exceed starting_luck value", err)
	case isFinancesValidationError(err):
		return badRequestError("character.invalid_finances", "credit_rating_skill_id must reference a skill from this character", err)
	case isDerivedStatsValidationError(err):
		return badRequestError("character.invalid_derived_stats", "derived stats payload is invalid", err)
	case isBackstoryItemValidationError(err):
		return badRequestError("character.invalid_backstory_section", "backstory item section is invalid", err)
	case isSkillValidationError(err):
		return badRequestError("character.invalid_skill", "skill payload is invalid", err)
	case isSkillReferencedError(err):
		return &myErrors.AppError{
			Status:  http.StatusConflict,
			Code:    "character.skill_in_use",
			Message: "skill is referenced by character finances",
			Err:     err,
		}
	default:
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: fallbackMessage,
			Err:     err,
		}
	}
}

func PortraitManagedByServerError() *myErrors.AppError {
	return badRequestError("character.portrait_managed_by_server", "portrait_url is managed by the server", nil)
}

func MapNotFoundOrServiceError(err error, notFoundMessage, fallbackMessage string) *myErrors.AppError {
	if errors.Is(err, pgx.ErrNoRows) {
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Code:    "character.not_found",
			Message: notFoundMessage,
			Err:     err,
		}
	}

	return MapServiceError(err, fallbackMessage)
}

func InvalidCharacterIDError(err error) *myErrors.AppError {
	return InvalidPathIDError("invalid character id", err)
}

func InvalidPathIDError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "character.invalid_id",
		Message: message,
		Err:     err,
	}
}

func InvalidInputError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "character.invalid_input",
		Message: message,
		Err:     err,
	}
}

func badRequestError(code string, message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func characterNameRequiredError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "character.name_required",
		Message: "name is required",
		Details: []myErrors.ErrorDetail{
			myErrors.ValidationDetail("body.name", "required"),
		},
		Err: err,
	}
}

func isHealthStateValidationError(err error) bool {
	return errors.Is(err, myErrors.ErrCurrentHealthExceedsMax)
}

func isMagicStateValidationError(err error) bool {
	return errors.Is(err, myErrors.ErrCurrentMagicExceedsMax)
}

func isSanityStateValidationError(err error) bool {
	return errors.Is(err, myErrors.ErrCurrentSanityExceedsMax)
}

func isLuckStateValidationError(err error) bool {
	return errors.Is(err, myErrors.ErrCurrentLuckExceedsStarting)
}

func isFinancesValidationError(err error) bool {
	return errors.Is(err, characterErrors.ErrInvalidFinances)
}

func isDerivedStatsValidationError(err error) bool {
	return errors.Is(err, characterErrors.ErrInvalidDerivedStats)
}

func isBackstoryItemValidationError(err error) bool {
	return errors.Is(err, characterErrors.ErrInvalidBackstorySection)
}

func isSkillValidationError(err error) bool {
	return errors.Is(err, characterErrors.ErrInvalidSkill)
}

func isSkillReferencedError(err error) bool {
	return errors.Is(err, characterErrors.ErrSkillInUse)
}
