package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// Characters
func (h *Handler) getAllCharacters(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characters, err := h.services.Character.GetAllCharacters(r.Context(), userID)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get user characters",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, characters)
	return nil
}

func (h *Handler) getCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	character, err := h.services.Character.GetCharacter(r.Context(), model.GetCharacterInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character data",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, character)
	return nil
}

func (h *Handler) createCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	var input db.CreateCharacterParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID

	character, err := h.services.Character.CreateCharacter(r.Context(), input)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create character",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusCreated, character)
	return nil
}

func (h *Handler) updateCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpdateCharacterParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.ID = characterID

	character, err := h.services.Character.UpdateCharacter(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update character",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, character)
	return nil
}

func (h *Handler) deleteCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteCharacter(r.Context(), db.DeleteCharacterParams{
		UserID: userID,
		ID:     characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Health
func (h *Handler) getHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	health, err := h.services.Character.GetHealth(r.Context(), db.GetHealthStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character health not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character health",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, health)
	return nil
}

func (h *Handler) upsertHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertHealthStateParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	health, err := h.services.Character.UpsertHealth(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}
		if isHealthStateValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "current_hp value cannot exceed max_hp value",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character health",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, health)
	return nil
}

func (h *Handler) deleteHealth(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteHealth(r.Context(), db.DeleteHealthStateParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character health not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character health",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Sanity
func (h *Handler) getSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	sanity, err := h.services.Character.GetSanity(r.Context(), db.GetSanityStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character sanity not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character sanity",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, sanity)
	return nil
}

func (h *Handler) upsertSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertSanityStateParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	sanity, err := h.services.Character.UpsertSanity(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		if isSanityStateValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "current_sanity value cannot exceed max_sanity value",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character sanity",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, sanity)
	return nil
}

func (h *Handler) deleteSanity(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteSanity(r.Context(), db.DeleteSanityStateParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character sanity not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character sanity",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Magic
func (h *Handler) getMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	magic, err := h.services.Character.GetMagic(r.Context(), db.GetMagicStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character magic not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character magic",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, magic)
	return nil
}

func (h *Handler) upsertMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertMagicStateParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	magic, err := h.services.Character.UpsertMagic(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		if isMagicStateValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "current_mp value cannot exceed max_mp value",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character magic",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, magic)
	return nil
}

func (h *Handler) deleteMagic(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteMagic(r.Context(), db.DeleteMagicStateParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character magic not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character magic",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Luck
func (h *Handler) getLuck(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	luck, err := h.services.Character.GetLuck(r.Context(), db.GetLuckStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character luck not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character luck",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, luck)
	return nil
}

func (h *Handler) upsertLuck(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertLuckStateParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	luck, err := h.services.Character.UpsertLuck(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		if isLuckStateValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "current_luck value cannot exceed starting_luck value",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character luck",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, luck)
	return nil
}

func (h *Handler) deleteLuck(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteLuck(r.Context(), db.DeleteLuckStateParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character luck not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character luck",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Finances
func (h *Handler) getFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	finances, err := h.services.Character.GetFinances(r.Context(), db.GetFinancesParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character finances not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character finances",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, finances)
	return nil
}

func (h *Handler) upsertFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertFinancesParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	finances, err := h.services.Character.UpsertFinances(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		if isFinancesValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "credit_rating_skill_id must reference a skill from this character",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character finances",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, finances)
	return nil
}

func (h *Handler) deleteFinances(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteFinances(r.Context(), db.DeleteFinancesParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character finances not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character finances",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Backstory
func (h *Handler) getBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	backstory, err := h.services.Character.GetBackstory(r.Context(), db.GetBackstoryByCharacterParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character backstory not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character backstory",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *Handler) upsertBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertBackstoryParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	backstory, err := h.services.Character.UpsertBackstory(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character backstory",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, backstory)
	return nil
}

func (h *Handler) deleteBackstory(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteBackstory(r.Context(), db.DeleteBackstoryParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character backstory not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character backstory",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) getBackstoryItems(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	items, err := h.services.Character.GetBackstoryItems(r.Context(), db.GetBackstoryItemsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character backstory items",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, items)
	return nil
}

func (h *Handler) getBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid backstory item id",
			Err:     err,
		}
	}

	item, err := h.services.Character.GetBackstoryItem(r.Context(), db.GetBackstoryItemParams{
		UserID:          userID,
		CharacterID:     characterID,
		BackstoryItemID: itemID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "backstory item not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get backstory item",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, item)
	return nil
}

func (h *Handler) createBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.CreateBackstoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	item, err := h.services.Character.CreateBackstoryItem(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character backstory not found",
				Err:     err,
			}
		}

		if isBackstoryItemValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "backstory item section is invalid",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create backstory item",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusCreated, item)
	return nil
}

func (h *Handler) updateBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid backstory item id",
			Err:     err,
		}
	}

	var input db.UpdateBackstoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.BackstoryItemID = itemID

	item, err := h.services.Character.UpdateBackstoryItem(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "backstory item not found",
				Err:     err,
			}
		}

		if isBackstoryItemValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "backstory item section is invalid",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update backstory item",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, item)
	return nil
}

func (h *Handler) deleteBackstoryItem(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	itemID, err := getBackstoryItemIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid backstory item id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteBackstoryItem(r.Context(), db.DeleteBackstoryItemParams{
		UserID:          userID,
		CharacterID:     characterID,
		BackstoryItemID: itemID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "backstory item not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete backstory item",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DerivedStats
func (h *Handler) getDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	derivedStats, err := h.services.Character.GetDerivedStats(r.Context(), db.GetDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character derived stats not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character derived stats",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, derivedStats)
	return nil
}

func (h *Handler) upsertDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertDerivedStatsParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	derivedStats, err := h.services.Character.UpsertDerivedStats(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		if isDerivedStatsValidationError(err) {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Message: "derived stats values cannot be negative",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character derived stats",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, derivedStats)
	return nil
}

func (h *Handler) deleteDerivedStats(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteDerivedStats(r.Context(), db.DeleteDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character derived stats not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character derived stats",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Characteristics
func (h *Handler) getCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	characteristics, err := h.services.Character.GetCharacteristics(r.Context(), db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get character characteristics",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, characteristics)
	return nil
}

func (h *Handler) upsertCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.UpsertCharacteristicsParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	characteristics, err := h.services.Character.UpsertCharacteristics(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to upsert character characteristics",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, characteristics)
	return nil
}

func (h *Handler) deleteCharacteristics(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteCharacteristics(r.Context(), db.DeleteCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character characteristics not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete character characteristics",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Notes
func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	notes, err := h.services.Character.GetNotes(r.Context(), db.GetNotesParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get all character notes",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, notes)
	return nil
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid note id",
			Err:     err,
		}
	}

	note, err := h.services.Character.GetNote(r.Context(), db.GetNoteParams{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "note not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get note data",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, note)
	return nil
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	var input db.CreateNoteParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID

	note, err := h.services.Character.CreateNote(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "character not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create note",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusCreated, note)
	return nil
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid note id",
			Err:     err,
		}
	}

	var input db.UpdateNoteParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.NoteID = noteID

	note, err := h.services.Character.UpdateNote(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "note not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update note",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, note)
	return nil
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid character id",
			Err:     err,
		}
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid note id",
			Err:     err,
		}
	}

	if err := h.services.Character.DeleteNote(r.Context(), db.DeleteNoteParams{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "note not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete note",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Utils
func getCharacterIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	characterID := chi.URLParam(r, "characterID")

	var characterUUID pgtype.UUID
	if err := characterUUID.Scan(characterID); err != nil {
		return pgtype.UUID{}, err
	}

	return characterUUID, nil
}

func getNoteIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	noteID := chi.URLParam(r, "noteID")

	var noteUUID pgtype.UUID
	if err := noteUUID.Scan(noteID); err != nil {
		return pgtype.UUID{}, err
	}

	return noteUUID, nil
}

func getBackstoryItemIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	itemID := chi.URLParam(r, "itemID")

	var itemUUID pgtype.UUID
	if err := itemUUID.Scan(itemID); err != nil {
		return pgtype.UUID{}, err
	}

	return itemUUID, nil
}

func isHealthStateValidationError(err error) bool {
	if errors.Is(err, myErrors.ErrCurrentHealthExceedsMax) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "chk_health_states_current_lte_max"
}

func isMagicStateValidationError(err error) bool {
	if errors.Is(err, myErrors.ErrCurrentMagicExceedsMax) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "chk_magic_states_current_lte_max"
}

func isSanityStateValidationError(err error) bool {
	if errors.Is(err, myErrors.ErrCurrentSanityExceedsMax) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "chk_sanity_states_current_lte_max"
}

func isLuckStateValidationError(err error) bool {
	if errors.Is(err, myErrors.ErrCurrentLuckExceedsStarting) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "chk_luck_states_current_lte_starting"
}

func isFinancesValidationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "fk_finances_credit_rating_skill"
}

func isDerivedStatsValidationError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	switch pgErr.ConstraintName {
	case "derived_stats_speed_check",
		"derived_stats_physique_check",
		"derived_stats_damage_bonus_check",
		"derived_stats_dodge_value_check":
		return true
	default:
		return false
	}
}

func isBackstoryItemValidationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.ConstraintName == "chk_backstory_items_section"
}
