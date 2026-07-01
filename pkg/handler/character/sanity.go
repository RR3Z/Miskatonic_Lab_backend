package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	sanity, err := h.service.GetSanity(r.Context(), sanityDTO.GetSanityInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character sanity not found", "failed to get character sanity")
	}

	utils.WriteJSON(w, http.StatusOK, sanity)
	return nil
}

func (h *CharacterHandler) upsertSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input sanityDTO.UpsertSanityInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	sanity, err := h.service.UpsertSanity(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character sanity")
	}

	utils.WriteJSON(w, http.StatusOK, sanity)
	return nil
}

func (h *CharacterHandler) deleteSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteSanity(r.Context(), sanityDTO.DeleteSanityInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character sanity not found", "failed to delete character sanity")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
