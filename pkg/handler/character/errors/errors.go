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

func MapNotFoundOrServiceError(err error, notFoundMessage, fallbackMessage string) *myErrors.AppError {
	if errors.Is(err, pgx.ErrNoRows) {
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
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
		Message: message,
		Err:     err,
	}
}

func InvalidInputError(message string, err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
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
