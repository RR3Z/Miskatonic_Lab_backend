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

// Characteristics
func (h *Handler) getCharacteristics(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_characteristics_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	characteristics, err := h.services.Character.GetCharacteristics(r.Context(), db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"characteristics not found",
				"component", "character_characteristics_api",
				"character_id", characterID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "character not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to get character characteristics",
			"component", "character_characteristics_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get character characteristics", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"successfully get character characteristics",
		"component", "character_characteristics_api",
		"character_id", characterID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, characteristics)
}

func (h *Handler) upsertCharacteristics(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_characteristics_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	var input db.UpsertCharacteristicsParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("invalid request body",
			"component", "character_characteristics_api",
			"error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input.UserID = userID
	input.CharacterID = characterID

	characteristics, err := h.services.Character.UpsertCharacteristics(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"character not found",
				"component", "character_characteristics_api",
				"user_id", userID,
				"character_id", characterID,
				"error", err,
			)
			http.Error(w, "character not found", http.StatusNotFound)
			return
		}

		slog.Error("failed to upsert character characteristics",
			"component", "character_characteristics_api",
			"user_id", userID,
			"character_id", characterID,
			"error", err)
		http.Error(w, "failed to upsert character characteristics", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"character characteristics upserted successfully",
		"component", "character_characteristics_api",
		"user_id", userID,
		"character_id", characterID,
	)

	utils.WriteJSON(w, http.StatusOK, characteristics)
}

func (h *Handler) deleteCharacteristics(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_characteristics_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	if err := h.services.Character.DeleteCharacteristics(r.Context(), db.DeleteCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"character characteristics not found",
				"component", "character_characteristics_api",
				"character_id", characterID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "character characteristics not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to delete character characteristics",
			"component", "character_characteristics_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to delete character characteristics", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"character characteristics deleted successfully",
		"component", "character_characteristics_api",
		"character_id", characterID,
		"user_id", userID,
	)

	w.WriteHeader(http.StatusNoContent)
}

// Notes
func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	notes, err := h.services.Character.GetNotes(r.Context(), db.GetNotesParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		slog.Error(
			"failed to get all character notes",
			"component", "character_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get all character notes", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"successfully get all character notes",
		"component", "character_note_api",
		"character_id", characterID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, notes)
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid note id format",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"error", err,
		)
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	note, err := h.services.Character.GetNote(r.Context(), db.GetNoteParams{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"note not found",
				"component", "character_note_api",
				"character_id", characterID,
				"note_id", noteID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to get note data",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get note data", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"successfully get note",
		"component", "character_note_api",
		"character_id", characterID,
		"note_id", noteID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, note)
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	var input db.CreateNoteParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error(
			"invalid request body",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input.UserID = userID
	input.CharacterID = characterID

	note, err := h.services.Character.CreateNote(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"character not found",
				"component", "character_note_api",
				"character_id", characterID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "character not found", http.StatusNotFound)
			return
		}

		slog.Error("failed to create note",
			"component", "character_note_api",
			"character_id", characterID,
			"user_id", userID,
			"error", err)
		http.Error(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"note created successfully",
		"component", "character_note_api",
		"character_id", characterID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusCreated, note)
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid note id format",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"error", err,
		)
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	var input db.UpdateNoteParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error(
			"invalid request body",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.NoteID = noteID

	note, err := h.services.Character.UpdateNote(r.Context(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"note not found",
				"component", "character_note_api",
				"character_id", characterID,
				"note_id", noteID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		slog.Error("failed to update note",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"user_id", userID,
			"error", err)
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"note updated successfully",
		"component", "character_note_api",
		"character_id", characterID,
		"note_id", noteID,
		"user_id", userID,
	)

	utils.WriteJSON(w, http.StatusOK, note)
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid character id format",
			"component", "character_note_api",
			"character_id", characterID,
			"error", err,
		)
		http.Error(w, "invalid character id", http.StatusBadRequest)
		return
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		slog.Error(
			"invalid note id format",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"error", err,
		)
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	if err := h.services.Character.DeleteNote(r.Context(), db.DeleteNoteParams{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(
				"note not found",
				"component", "character_note_api",
				"character_id", characterID,
				"note_id", noteID,
				"user_id", userID,
				"error", err,
			)
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"failed to delete note",
			"component", "character_note_api",
			"character_id", characterID,
			"note_id", noteID,
			"user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to delete note", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"note deleted successfully",
		"component", "character_api",
		"character_id", characterID,
		"note_id", noteID,
		"user_id", userID,
	)

	w.WriteHeader(http.StatusNoContent)
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
