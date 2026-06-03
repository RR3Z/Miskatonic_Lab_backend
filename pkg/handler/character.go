package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Characters
func (h *Handler) getAllCharacters(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characters, err := h.services.Character.GetAllCharacters(r.Context(), userID)
	if err != nil {
		slog.Error(
			"failed to get all user characters",
			"component", "character_api",
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get all user characters", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"successfully get all user characters",
		"component", "character_api",
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, characters)
}

func (h *Handler) getCharacter(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	character, err := h.services.Character.GetCharacter(r.Context(), model.GetCharacterInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"character not found",
				"component", "character_api",
				"character_id", characterID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "character not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to get character data",
			"component", "character_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get character data", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"successfully get character",
		"component", "character_api",
		"character_id", characterID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, character)
}

func (h *Handler) createCharacter(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	var input db.CreateCharacterParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("invalid request body", "component", "character_api", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input.UserID = userID

	character, err := h.services.Character.CreateCharacter(r.Context(), input)
	if err != nil {
		slog.Error("failed to create character",
			"component", "character_api",
			"user_id", userID,
			"error", err)
		http.Error(w, "failed to create character", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"character created successfully",
		"component", "character_api",
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusCreated, character)
}

func (h *Handler) updateCharacter(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	var input db.UpdateCharacterParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("invalid request body",
			"component", "character_api",
			"error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input.UserID = userID
	input.ID = characterID

	character, err := h.services.Character.UpdateCharacter(r.Context(), input)
	if err != nil {
		slog.Error("failed to update character",
			"component", "character_api",
			"user_id", userID,
			"error", err)
		http.Error(w, "failed to update character", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"character updated successfully",
		"component", "character_api",
		"character_id", characterID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, character)

}

func (h *Handler) deleteCharacter(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	if err := h.services.Character.DeleteCharacter(r.Context(), db.DeleteCharacterParams{
		UserID: userID,
		ID:     characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"character not found",
				"component", "character_api",
				"character_id", characterID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "character not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to delete character",
			"component", "character_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to delete character", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"character deleted successfully",
		"component", "character_api",
		"character_id", characterID,
		"user_id", userID,
	)

	w.WriteHeader(http.StatusNoContent)
}

func getCharacterIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	characterID := chi.URLParam(r, "characterID")

	var characterUUID pgtype.UUID
	if err := characterUUID.Scan(characterID); err != nil {
		return pgtype.UUID{}, err
	}

	return characterUUID, nil
}
