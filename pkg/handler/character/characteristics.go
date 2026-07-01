package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	characteristics, err := h.service.GetCharacteristics(r.Context(), characteristicsDTO.GetCharacteristicsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to get character characteristics")
	}

	utils.WriteJSON(w, http.StatusOK, characteristics)
	return nil
}

func (h *CharacterHandler) upsertCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input characteristicsDTO.UpsertCharacteristicsInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	characteristics, err := h.service.UpsertCharacteristics(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character characteristics")
	}

	utils.WriteJSON(w, http.StatusOK, characteristics)
	return nil
}

func (h *CharacterHandler) deleteCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteCharacteristics(r.Context(), characteristicsDTO.DeleteCharacteristicsInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character characteristics not found", "failed to delete character characteristics")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
