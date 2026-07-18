package characterHelpers

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func DecodeJSON(r *http.Request, target any) *myErrors.AppError {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return characterErrors.InvalidInputError("invalid request body", err)
	}
	return nil
}

func GetCharacterIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "characterID")
}

func GetNoteIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "noteID")
}

func GetInventoryItemIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "itemID")
}

func GetBackstoryItemIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "itemID")
}

func GetSkillIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "skillID")
}

func getUUIDFromRequest(r *http.Request, param string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(chi.URLParam(r, param)); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}
