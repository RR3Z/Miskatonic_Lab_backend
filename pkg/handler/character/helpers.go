package character

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func decodeJSON(r *http.Request, target any) *myErrors.AppError {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return characterErrors.InvalidInputError("invalid request body", err)
	}
	return nil
}

func getUUIDFromRequest(r *http.Request, param string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(chi.URLParam(r, param)); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}

func getCharacterIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "characterID")
}

func getNoteIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "noteID")
}

func getBackstoryItemIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "itemID")
}

func getSkillIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	return getUUIDFromRequest(r, "skillID")
}
