package diceRollerHelpers

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	diceRollerErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller/errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetCharacterIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(chi.URLParam(r, "characterID")); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}

func DecodeJSON(r *http.Request, target any) *myErrors.AppError {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return diceRollerErrors.InvalidInputError("invalid request body", err)
	}
	return nil
}
