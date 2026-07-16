package diceRollerDTO

import (
	"encoding/json"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type DiceRollModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	UserID      string             `json:"user_id"`
	Expression  string             `json:"expression"`
	Result      int32              `json:"result"`
	Details     json.RawMessage    `json:"details"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
}

func ToDiceRollModel(r db.DiceRoll) DiceRollModel {
	return DiceRollModel{
		ID:          r.ID,
		CharacterID: r.CharacterID,
		UserID:      r.UserID,
		Expression:  r.Expression,
		Result:      r.Result,
		Details:     json.RawMessage(r.Details),
		CreatedAt:   r.CreatedAt,
	}
}
