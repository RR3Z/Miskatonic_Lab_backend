package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getSkills(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	skills, err := h.service.GetSkills(r.Context(), skillsDTO.GetSkillsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get character skills")
	}

	utils.WriteJSON(w, http.StatusOK, skills)
	return nil
}

func (h *CharacterHandler) getSkill(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	skillID, err := characterHelpers.GetSkillIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid skill id", err)
	}

	skill, err := h.service.GetSkill(r.Context(), skillsDTO.GetSkillInput{
		UserID:      userID,
		CharacterID: characterID,
		SkillID:     skillID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "skill not found", "failed to get skill")
	}

	utils.WriteJSON(w, http.StatusOK, skill)
	return nil
}

func (h *CharacterHandler) createSkill(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input skillsDTO.CreateSkillInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	skill, err := h.service.CreateSkill(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to create skill")
	}

	utils.WriteJSON(w, http.StatusCreated, skill)
	return nil
}

func (h *CharacterHandler) updateSkill(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	skillID, err := characterHelpers.GetSkillIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid skill id", err)
	}

	var input skillsDTO.UpdateSkillInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.SkillID = skillID

	skill, err := h.service.UpdateSkill(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "skill not found", "failed to update skill")
	}

	utils.WriteJSON(w, http.StatusOK, skill)
	return nil
}

func (h *CharacterHandler) deleteSkill(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	skillID, err := characterHelpers.GetSkillIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid skill id", err)
	}

	if err := h.service.DeleteSkill(r.Context(), skillsDTO.DeleteSkillInput{
		UserID:      userID,
		CharacterID: characterID,
		SkillID:     skillID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "skill not found", "failed to delete skill")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
