package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getNotes(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	notes, err := h.service.GetNotes(r.Context(), notesDTO.GetNotesInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to get all character notes")
	}

	utils.WriteJSON(w, http.StatusOK, notes)
	return nil
}

func (h *CharacterHandler) getNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := characterHelpers.GetNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	note, err := h.service.GetNote(r.Context(), notesDTO.GetNoteInput{
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

func (h *CharacterHandler) createNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	var req notesDTO.NoteRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	input := notesDTO.CreateNoteInput{
		Title:       req.Title,
		Body:        req.Body,
		UserID:      userID,
		CharacterID: characterID,
	}

	note, err := h.service.CreateNote(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to create note")
	}

	utils.WriteJSON(w, http.StatusCreated, note)
	return nil
}

func (h *CharacterHandler) updateNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := characterHelpers.GetNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	var req notesDTO.NoteRequest
	if appErr := characterHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}
	input := notesDTO.UpdateNoteInput{
		Title:       req.Title,
		Body:        req.Body,
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	}

	note, err := h.service.UpdateNote(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "note not found", "failed to update note")
	}

	utils.WriteJSON(w, http.StatusOK, note)
	return nil
}

func (h *CharacterHandler) deleteNote(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	noteID, err := characterHelpers.GetNoteIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidPathIDError("invalid note id", err)
	}

	if err := h.service.DeleteNote(r.Context(), notesDTO.DeleteNoteInput{
		UserID:      userID,
		CharacterID: characterID,
		NoteID:      noteID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "note not found", "failed to delete note")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
