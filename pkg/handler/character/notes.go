package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	notes, err := h.characters.GetNotes(r.Context(), model.GetNotesInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get all character notes")
	}

	utils.WriteJSON(w, http.StatusOK, notes)
	return nil
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	note, err := h.characters.GetNote(r.Context(), model.GetNoteInput{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "note not found", "failed to get note data")
	}

	utils.WriteJSON(w, http.StatusOK, note)
	return nil
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var input model.CreateNoteInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID

	note, err := h.characters.CreateNote(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to create note")
	}

	utils.WriteJSON(w, http.StatusCreated, note)
	return nil
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	var input model.UpdateNoteInput
	if appErr := decodeJSON(r, &input); appErr != nil {
		return appErr
	}
	input.UserID = userID
	input.CharacterID = characterID
	input.NoteID = noteID

	note, err := h.characters.UpdateNote(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "note not found", "failed to update note")
	}

	utils.WriteJSON(w, http.StatusOK, note)
	return nil
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := getCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := getNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	if err := h.characters.DeleteNote(r.Context(), model.DeleteNoteInput{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "note not found", "failed to delete note")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
