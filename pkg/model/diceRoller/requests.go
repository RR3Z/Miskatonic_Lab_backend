package diceRollerDTO

import "github.com/jackc/pgx/v5/pgtype"

type MakeRollRequest struct {
	Expression string       `json:"expression"`
	RoomID     *pgtype.UUID `json:"room_id,omitempty"`
}
