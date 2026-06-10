package diceRoller

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type IDiceRoller interface {
	GetLastDiceRolls(ctx context.Context, input db.GetDiceRollsParams) ([]db.DiceRoll, error)
	MakeRoll(ctx context.Context, input DiceRollInput) (db.DiceRoll, error)
}

type DiceRollInput struct {
	UserID      string
	CharacterID pgtype.UUID
	Formula     string
}
