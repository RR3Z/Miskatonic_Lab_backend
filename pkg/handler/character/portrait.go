package character

import (
	"errors"
	"io"
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

func (h *CharacterHandler) updatePortrait(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
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
			return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitRequired, "failed to update character portrait")
		}
		if nextErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(nextErr, &maxBytesError) {
				return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitTooLarge, "failed to update character portrait")
			}
			return characterErrors.InvalidInputError("invalid portrait upload", nextErr)
		}

		if part.FormName() != "portrait" || part.FileName() == "" {
			_ = part.Close()
			continue
		}

		character, serviceErr := h.service.UpdatePortrait(r.Context(), characterDTO.UpdateCharacterPortraitInput{
			UserID:      utils.GetUserIDFromContext(r.Context()),
			CharacterID: characterID,
			Content:     part,
		})
		_ = part.Close()
		if serviceErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(serviceErr, &maxBytesError) {
				return characterErrors.MapServiceError(characterServiceErrors.ErrPortraitTooLarge, "failed to update character portrait")
			}
			return characterErrors.MapNotFoundOrServiceError(serviceErr, "character not found", "failed to update character portrait")
		}

		utils.WriteJSON(w, http.StatusOK, character)
		return nil
	}
}
