package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	derivedStats, err := h.service.GetDerivedStats(r.Context(), model.GetDerivedStatsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character derived stats not found", "failed to get character derived stats")
	}

	utils.WriteJSON(w, http.StatusOK, derivedStats)
	return nil
}

func (h *CharacterHandler) upsertDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.UpsertDerivedStatsInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	derivedStats, err := h.service.UpsertDerivedStats(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character derived stats")
	}

	utils.WriteJSON(w, http.StatusOK, derivedStats)
	return nil
}

func (h *CharacterHandler) deleteDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteDerivedStats(r.Context(), model.DeleteDerivedStatsInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character derived stats not found", "failed to delete character derived stats")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
