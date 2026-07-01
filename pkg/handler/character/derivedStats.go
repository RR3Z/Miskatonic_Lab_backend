package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	derivedStats, err := h.service.GetDerivedStats(r.Context(), derivedStatsDTO.GetDerivedStatsInput{
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

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var req derivedStatsDTO.DerivedStatsRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	input := derivedStatsDTO.UpsertDerivedStatsInput{
		UserID:      userID,
		CharacterID: characterID,
		Speed:       req.Speed,
		Physique:    req.Physique,
		DamageBonus: req.DamageBonus,
		DodgeValue:  req.DodgeValue,
	}

	derivedStats, err := h.service.UpsertDerivedStats(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to upsert character derived stats")
	}

	utils.WriteJSON(w, http.StatusOK, derivedStats)
	return nil
}

func (h *CharacterHandler) deleteDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteDerivedStats(r.Context(), derivedStatsDTO.DeleteDerivedStatsInput{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character derived stats not found", "failed to delete character derived stats")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
