package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	finances, err := h.service.GetFinances(r.Context(), financesDTO.GetFinancesInput{
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

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var req financesDTO.FinancesRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	input := financesDTO.UpsertFinancesInput{
		UserID:              userID,
		CharacterID:         characterID,
		SpendingLimit:       req.SpendingLimit,
		Cash:                req.Cash,
		Assets:              req.Assets,
		CreditRatingSkillID: req.CreditRatingSkillID,
	}

	finances, err := h.service.UpsertFinances(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character finances")
	}

	utils.WriteJSON(w, http.StatusOK, finances)
	return nil
}

func (h *CharacterHandler) deleteFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteFinances(r.Context(), financesDTO.DeleteFinancesInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character finances not found", "failed to delete character finances")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
