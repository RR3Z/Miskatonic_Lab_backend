package diceRollerDTO

import "github.com/jackc/pgx/v5/pgtype"

type MakeRollRequest struct {
	Expression string       `json:"expression"`
	D100Mode   *D100Mode    `json:"d100_mode,omitempty"`
	RoomID     *pgtype.UUID `json:"room_id,omitempty"`
}
