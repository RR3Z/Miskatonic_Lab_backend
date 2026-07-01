package roomHelpers

import (
	"net/http"

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
