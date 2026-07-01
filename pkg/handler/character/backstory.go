package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *Handler) getBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	backstory, err := h.characters.GetBackstory(r.Context(), model.GetBackstoryInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to get character backstory")
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *Handler) upsertBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.UpsertBackstoryInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	backstory, err := h.characters.UpsertBackstory(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character backstory")
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *Handler) deleteBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.characters.DeleteBackstory(r.Context(), model.DeleteBackstoryInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to delete character backstory")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) getBackstoryItems(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	items, err := h.characters.GetBackstoryItems(r.Context(), model.GetBackstoryItemsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get character backstory items")
	}

	utils.WriteJSON(w, http.StatusOK, items)
	return nil
}

func (h *Handler) getBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	item, err := h.characters.GetBackstoryItem(r.Context(), model.GetBackstoryItemInput{
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

func (h *Handler) createBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.CreateBackstoryItemInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	item, err := h.characters.CreateBackstoryItem(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character backstory not found", "failed to create backstory item")
	}

	utils.WriteJSON(w, http.StatusCreated, item)
	return nil
}

func (h *Handler) updateBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	var input model.UpdateBackstoryItemInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.BackstoryItemID = itemID

	item, err := h.characters.UpdateBackstoryItem(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "backstory item not found", "failed to update backstory item")
	}

	utils.WriteJSON(w, http.StatusOK, item)
	return nil
}

func (h *Handler) deleteBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid backstory item id", err)
	}

	if err := h.characters.DeleteBackstoryItem(r.Context(), model.DeleteBackstoryItemInput{
		UserID:          userID,
		CharacterID:     characterID,
		BackstoryItemID: itemID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "backstory item not found", "failed to delete backstory item")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
