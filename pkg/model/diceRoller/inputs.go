package diceRollerDTO

import "github.com/jackc/pgx/v5/pgtype"

type MakeRollInput struct {
	UserID      string
	CharacterID pgtype.UUID
	Formula     string
	D100Mode    *D100Mode
	RoomID      *pgtype.UUID
}

type GetLastDiceRollsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
