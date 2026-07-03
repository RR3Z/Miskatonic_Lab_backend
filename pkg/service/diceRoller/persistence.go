package diceRoller

import (
	"context"
	"errors"
	"log/slog"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
)

func (s *DiceRollerService) makeLocalRoll(
	ctx context.Context,
	input diceRollerDTO.MakeRollInput,
	prepared preparedDiceRoll,
) (diceRollerDTO.DiceRollModel, error) {
	diceRoll, err := createDiceRoll(ctx, s.repos.Queries, input, prepared)
	if err != nil {
		return diceRollerDTO.DiceRollModel{}, err
	}

	cleanOldDiceRolls(ctx, s.repos.Queries, input)

	return diceRollerDTO.ToDiceRollModel(diceRoll), nil
}

func createDiceRoll(
	ctx context.Context,
	queries *db.Queries,
	input diceRollerDTO.MakeRollInput,
	prepared preparedDiceRoll,
) (db.DiceRoll, error) {
	diceRoll, err := queries.CreateDiceRoll(ctx, db.CreateDiceRollParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		Expression:  input.Formula,
		Result:      int32(prepared.result),
		Details:     prepared.detailsJSON,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.DiceRoll{}, ErrCharacterNotFound
		}
		return db.DiceRoll{}, err
	}

	return diceRoll, nil
}

func cleanOldDiceRolls(ctx context.Context, queries *db.Queries, input diceRollerDTO.MakeRollInput) {
	if err := queries.CleanOldDiceRolls(ctx, db.CleanOldDiceRollsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		slog.Warn("failed to clean old dice rolls", "character_id", input.CharacterID, "error", err)
	}
}
