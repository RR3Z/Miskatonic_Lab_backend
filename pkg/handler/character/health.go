package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	health, err := h.service.GetHealth(r.Context(), model.GetHealthInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character health not found", "failed to get character health")
	}

	utils.WriteJSON(w, http.StatusOK, health)
	return nil
}

func (h *CharacterHandler) upsertHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.UpsertHealthInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	health, err := h.service.UpsertHealth(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character health")
	}

	utils.WriteJSON(w, http.StatusOK, health)
	return nil
}

func (h *CharacterHandler) deleteHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteHealth(r.Context(), model.DeleteHealthInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character health not found", "failed to delete character health")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
