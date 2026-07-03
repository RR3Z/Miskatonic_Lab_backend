package roomDTO

import "github.com/jackc/pgx/v5/pgtype"

type LeaveRoomResult struct {
	DeletedRoomID *pgtype.UUID `json:"deleted_room_id,omitempty"`
}
