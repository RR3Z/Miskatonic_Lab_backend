package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *CharacterHandler) getInventoryItems(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	items, err := h.service.GetInventoryItems(r.Context(), inventoryDTO.GetInventoryItemsInput{UserID: userID, CharacterID: characterID})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get character inventory")
	}
	utils.WriteJSON(w, http.StatusOK, inventoryDTO.ToInventoryItemModels(items))
	return nil
}

func (h *CharacterHandler) getInventoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	characterID, itemID, appErr := inventoryItemIDs(r)
	if appErr != nil {
		return appErr
	}

	item, err := h.service.GetInventoryItem(r.Context(), inventoryDTO.GetInventoryItemInput{UserID: userID, CharacterID: characterID, ItemID: itemID})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "inventory item not found", "failed to get inventory item")
	}
	utils.WriteJSON(w, http.StatusOK, inventoryDTO.ToInventoryItemModel(item))
	return nil
}

func (h *CharacterHandler) createInventoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var req inventoryDTO.InventoryItemRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	item, err := h.service.CreateInventoryItem(r.Context(), inventoryDTO.CreateInventoryItemInput{
		Name:        req.Name,
		Quantity:    req.Quantity,
		Category:    req.Category,
		Description: req.Description,
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to create inventory item")
	}
	utils.WriteJSON(w, http.StatusCreated, inventoryDTO.ToInventoryItemModel(item))
	return nil
}

func (h *CharacterHandler) updateInventoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	characterID, itemID, appErr := inventoryItemIDs(r)
	if appErr != nil {
		return appErr
	}

	var req inventoryDTO.InventoryItemRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	item, err := h.service.UpdateInventoryItem(r.Context(), inventoryDTO.UpdateInventoryItemInput{
		Name:        req.Name,
		Quantity:    req.Quantity,
		Category:    req.Category,
		Description: req.Description,
		UserID:      userID,
		CharacterID: characterID,
		ItemID:      itemID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "inventory item not found", "failed to update inventory item")
	}
	utils.WriteJSON(w, http.StatusOK, inventoryDTO.ToInventoryItemModel(item))
	return nil
}

func (h *CharacterHandler) deleteInventoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	characterID, itemID, appErr := inventoryItemIDs(r)
	if appErr != nil {
		return appErr
	}

	if err := h.service.DeleteInventoryItem(r.Context(), inventoryDTO.DeleteInventoryItemInput{
		UserID:      userID,
		CharacterID: characterID,
		ItemID:      itemID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "inventory item not found", "failed to delete inventory item")
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func inventoryItemIDs(r *http.Request) (pgtype.UUID, pgtype.UUID, *myErrors.AppError) {
	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, characterErrors.InvalidCharacterIDError(err)
	}

	itemID, err := characterHelpers.GetInventoryItemIDFromRequest(r)
	if err != nil {
		return pgtype.UUID{}, pgtype.UUID{}, characterErrors.InvalidPathIDError("invalid inventory item id", err)
	}

	return characterID, itemID, nil
}
