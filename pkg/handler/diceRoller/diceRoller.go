package diceRoller

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	diceRollerErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller/errors"
	diceRollerHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller/helpers"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *DiceRollerHandler) makeRoll(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := diceRollerHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return diceRollerErrors.InvalidCharacterIDError(err)
	}

	var req diceRollerDTO.MakeRollRequest
	if appErr := diceRollerHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	if req.RoomID != nil && req.RoomID.Valid {
		if h.roomChecker == nil {
			return diceRollerErrors.RoomNotAvailableError(nil, "room dice rolls are not available")
		}
		if checkErr := h.roomChecker.EnsureCanPublishRoomEvent(r.Context(), *req.RoomID, userID); checkErr != nil {
			return diceRollerErrors.RoomNotAvailableError(checkErr, "room not available for dice roll")
		}
	}

	roll, err := h.service.MakeRoll(r.Context(), diceRollerDTO.MakeRollInput{
		UserID:      userID,
		CharacterID: characterID,
		Formula:     req.Expression,
		RoomID:      req.RoomID,
	})
	if err != nil {
		return diceRollerErrors.MapServiceError(err, "failed to create dice roll")
	}

	utils.WriteJSON(w, http.StatusCreated, roll)
	return nil
}

func (h *DiceRollerHandler) getLastDiceRolls(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := diceRollerHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return diceRollerErrors.InvalidCharacterIDError(err)
	}

	rolls, err := h.service.GetLastDiceRolls(r.Context(), diceRollerDTO.GetLastDiceRollsInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return diceRollerErrors.MapServiceError(err, "failed to get dice rolls")
	}

	utils.WriteJSON(w, http.StatusOK, rolls)
	return nil
}
