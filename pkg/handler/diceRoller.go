package handler

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *Handler) getLastDiceRolls(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	rolls, err := h.services.DiceRoller.GetLastDiceRolls(r.Context(), db.GetDiceRollsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get dice rolls",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, rolls)
	return nil
}

func (h *Handler) makeRoll(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var req model.MakeRollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	if req.Expression == "" {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "expression is required",
		}
	}

	roll, err := h.services.DiceRoller.MakeRoll(r.Context(), diceRoller.DiceRollInput{
		UserID:      userID,
		CharacterID: characterID,
		Formula:     req.Expression,
	})
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create dice roll",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusCreated, roll)
	return nil
}
