package character

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/helpers"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
)

type characterWriteRequestEnvelope struct {
	characterDTO.CharacterRequest
	ForbiddenPortraitURL json.RawMessage `json:"portrait_url"`
}

func decodeCharacterPatchRequest(r *http.Request) (characterDTO.PatchCharacterRequest, *myErrors.AppError) {
	var request characterDTO.PatchCharacterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return characterDTO.PatchCharacterRequest{}, characterErrors.InvalidInputError("invalid request body", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return characterDTO.PatchCharacterRequest{}, characterErrors.InvalidInputError("invalid request body", err)
	}
	if !request.HasChanges() || request.Name.Set && request.Name.Value == nil {
		return characterDTO.PatchCharacterRequest{}, characterErrors.InvalidInputError("character patch must contain at least one valid profile field", nil)
	}

	return request, nil
}

func decodeCharacterWriteRequest(r *http.Request) (characterDTO.CharacterRequest, *myErrors.AppError) {
	var envelope characterWriteRequestEnvelope
	if appErr := characterHelpers.DecodeJSON(r, &envelope); appErr != nil {
		return characterDTO.CharacterRequest{}, appErr
	}
	if len(envelope.ForbiddenPortraitURL) > 0 {
		return characterDTO.CharacterRequest{}, characterErrors.PortraitManagedByServerError()
	}
	return envelope.CharacterRequest, nil
}
