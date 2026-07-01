package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	magic, err := h.service.GetMagic(r.Context(), magicDTO.GetMagicInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character magic not found", "failed to get character magic")
	}

	utils.WriteJSON(w, http.StatusOK, magic)
	return nil
}

func (h *CharacterHandler) upsertMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input magicDTO.UpsertMagicInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	magic, err := h.service.UpsertMagic(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character magic")
	}

	utils.WriteJSON(w, http.StatusOK, magic)
	return nil
}

func (h *CharacterHandler) deleteMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteMagic(r.Context(), magicDTO.DeleteMagicInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character magic not found", "failed to delete character magic")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
