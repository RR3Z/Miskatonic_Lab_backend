package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	finances, err := h.service.GetFinances(r.Context(), model.GetFinancesInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character finances not found", "failed to get character finances")
	}

	utils.WriteJSON(w, http.StatusOK, finances)
	return nil
}

func (h *CharacterHandler) upsertFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.UpsertFinancesInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	finances, err := h.service.UpsertFinances(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character finances")
	}

	utils.WriteJSON(w, http.StatusOK, finances)
	return nil
}

func (h *CharacterHandler) deleteFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteFinances(r.Context(), model.DeleteFinancesInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character finances not found", "failed to delete character finances")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
