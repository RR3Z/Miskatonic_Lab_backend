package diceRollerHelpers

import "github.com/jackc/pgx/v5/pgtype"

func EventRoomID(roomID *pgtype.UUID) *string {
	if roomID == nil || !roomID.Valid {
		return nil
	}

	value := roomID.String()
	return &value
}
