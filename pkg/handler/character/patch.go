package character

import (
	"errors"
	"io"
	"mime"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	characterServiceErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
)

const MaxPortraitUploadBytes = portraitStorage.MaxUploadBytes

func (h *CharacterHandler) patchCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return characterErrors.InvalidInputError("invalid content type", err)
	}

	switch mediaType {
	case "multipart/form-data":
		return h.replacePortrait(w, r)
	case "application/json":
		return h.patchCharacterProfile(w, r)
	default:
		return characterErrors.InvalidInputError("content type must be application/json or multipart/form-data", nil)
	}
}

func (h *CharacterHandler) patchCharacterProfile(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	request, appErr := decodeCharacterPatchRequest(r)
	if appErr != nil {
		return appErr
	}

	character, err := h.service.PatchCharacter(r.Context(), characterDTO.PatchCharacterInput{
		UserID:     utils.GetUserIDFromContext(r.Context()),
		ID:         characterID,
		Name:       request.Name,
		Occupation: request.Occupation,
		Age:        request.Age,
		Sex:        request.Sex,
		Residence:  request.Residence,
		Birthplace: request.Birthplace,
	})
	if err != nil {
		return characterErrors.MapNotFoundOrServiceError(err, "character not found", "failed to patch character")
	}

	utils.WriteJSON(w, http.StatusOK, character)
	return nil
}

func (h *CharacterHandler) replacePortrait(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	characterID, err := characterHelpers.GetCharacterIDFromRequest(r)
	if err != nil {
		return characterErrors.InvalidCharacterIDError(err)
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxPortraitUploadBytes+(1<<20))
	reader, err := r.MultipartReader()
	if err != nil {
		return characterErrors.InvalidInputError("invalid portrait upload", err)
	}

	for {
		part, nextErr := reader.NextPart()
		if errors.Is(nextErr, io.EOF) {
			return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitRequired, "failed to replace character portrait")
		}
		if nextErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(nextErr, &maxBytesError) {
				return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitTooLarge, "failed to replace character portrait")
			}
			return characterErrors.InvalidInputError("invalid portrait upload", nextErr)
		}

		if part.FormName() != "portrait" || part.FileName() == "" {
			_ = part.Close()
			continue
		}

		character, serviceErr := h.service.ReplacePortrait(r.Context(), characterDTO.ReplacePortraitInput{
			UserID:      utils.GetUserIDFromContext(r.Context()),
			CharacterID: characterID,
			File:        part,
		})
		_ = part.Close()
		if serviceErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(serviceErr, &maxBytesError) {
				return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitTooLarge, "failed to replace character portrait")
			}
			return characterErrors.MapNotFoundOrServiceError(serviceErr, "character not found", "failed to replace character portrait")
		}

		utils.WriteJSON(w, http.StatusOK, character)
		return nil
	}
}
