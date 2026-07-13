package character

import (
	"encoding/json"
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
