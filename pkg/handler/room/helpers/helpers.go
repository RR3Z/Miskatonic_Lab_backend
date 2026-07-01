package roomHelpers

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	roomErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room/errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetRoomIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(chi.URLParam(r, "roomID")); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}

func DecodeJSON(r *http.Request, target any) *myErrors.AppError {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return roomErrors.InvalidInputError("invalid request body", err)
	}
	return nil
}
