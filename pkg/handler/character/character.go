package character

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

func (h *CharacterHandler) getAllCharacters(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characters, err := h.service.GetAllCharacters(r.Context(), userID)
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

func (h *CharacterHandler) getCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	character, err := h.service.GetCharacter(r.Context(), characterDTO.GetCharacterInput{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to get character data")
	}

	utils.WriteJSON(w, http.StatusOK, character)
	return nil
}

func (h *CharacterHandler) createCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	req, appErr := decodeCharacterWriteRequest(r)
	if appErr != nil {
		return appErr
	}
	input := characterDTO.CreateCharacterInput{
		UserID:     userID,
		Name:       req.Name,
		PlayerName: req.PlayerName,
		Occupation: req.Occupation,
		Age:        req.Age,
		Sex:        req.Sex,
		Residence:  req.Residence,
		Birthplace: req.Birthplace,
	}

	character, err := h.service.CreateCharacter(r.Context(), input)
	if err != nil {
		return characterErrors.MapServiceError(err, "failed to create character")
	}

	utils.WriteJSON(w, http.StatusCreated, character)
	return nil
}

func (h *CharacterHandler) updateCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	req, appErr := decodeCharacterWriteRequest(r)
	if appErr != nil {
		return appErr
	}
	input := characterDTO.UpdateCharacterInput{
		UserID:     userID,
		ID:         characterID,
		Name:       req.Name,
		PlayerName: req.PlayerName,
		Occupation: req.Occupation,
		Age:        req.Age,
		Sex:        req.Sex,
		Residence:  req.Residence,
		Birthplace: req.Birthplace,
	}

	character, err := h.service.UpdateCharacter(r.Context(), input)
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to update character")
	}

	utils.WriteJSON(w, http.StatusOK, character)
	return nil
}

func (h *CharacterHandler) deleteCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	if err := h.service.DeleteCharacter(r.Context(), characterDTO.DeleteCharacterInput{
		UserID: userID,
		ID:     characterID,
	}); err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to delete character")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
