package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	backstory, err := h.service.GetBackstory(r.Context(), backstoriesDTO.GetBackstoryInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to get character backstory")
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *CharacterHandler) upsertBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input backstoriesDTO.UpsertBackstoryInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	backstory, err := h.service.UpsertBackstory(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character backstory")
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *CharacterHandler) deleteBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteBackstory(r.Context(), backstoriesDTO.DeleteBackstoryInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to delete character backstory")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CharacterHandler) getBackstoryItems(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	items, err := h.service.GetBackstoryItems(r.Context(), backstoriesDTO.GetBackstoryItemsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get character backstory items")
	}

	utils.WriteJSON(w, http.StatusOK, items)
	return nil
}

func (h *CharacterHandler) getBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := characterHelpers.GetBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	item, err := h.service.GetBackstoryItem(r.Context(), backstoriesDTO.GetBackstoryItemInput{
		UserID:          userID,
		CharacterID:     characterID,
		BackstoryItemID: itemID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "backstory item not found", "failed to get backstory item")
	}

	utils.WriteJSON(w, http.StatusOK, item)
	return nil
}

func (h *CharacterHandler) createBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input backstoriesDTO.CreateBackstoryItemInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	item, err := h.service.CreateBackstoryItem(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to create backstory item")
	}

	utils.WriteJSON(w, http.StatusCreated, item)
	return nil
}

func (h *CharacterHandler) updateBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := characterHelpers.GetBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	var input backstoriesDTO.UpdateBackstoryItemInput
	if appErr := characterHelpers.DecodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.BackstoryItemID = itemID

	item, err := h.service.UpdateBackstoryItem(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "backstory item not found", "failed to update backstory item")
	}

	utils.WriteJSON(w, http.StatusOK, item)
	return nil
}

func (h *CharacterHandler) deleteBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := characterHelpers.GetBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	if err := h.service.DeleteBackstoryItem(r.Context(), backstoriesDTO.DeleteBackstoryItemInput{
		UserID:          userID,
		CharacterID:     characterID,
		BackstoryItemID: itemID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "backstory item not found", "failed to delete backstory item")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
