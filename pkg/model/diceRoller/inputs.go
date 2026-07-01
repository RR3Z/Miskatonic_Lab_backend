package diceRollerDTO

import "github.com/jackc/pgx/v5/pgtype"

type MakeRollInput struct {
	UserID      string
	CharacterID pgtype.UUID
	Formula     string
}

type GetLastDiceRollsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
