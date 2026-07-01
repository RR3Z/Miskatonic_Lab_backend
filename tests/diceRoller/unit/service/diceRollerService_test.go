package tests

import (
	"context"
	"testing"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func diceRollServiceTestUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}

func TestMakeRoll_ValidationFailsBeforePersistence_Empty(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFailsBeforePersistence_NoDice(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "5",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFailsBeforePersistence_Subtraction(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "3d6-1d4",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFailsBeforePersistence_ZeroCount(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "0d6",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFailsBeforePersistence_ZeroSides(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "1d0",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFailsBeforePersistence_NoCountBeforeD(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "d20",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}

func TestMakeRoll_ValidationFails_BadSyntax(t *testing.T) {
	service := diceRollerServices.NewDiceRollerService(&repository.Repository{})

	_, err := service.MakeRoll(context.Background(), diceRollerDTO.MakeRollInput{
		UserID:      "user_1",
		CharacterID: diceRollServiceTestUUID("11111111-1111-1111-1111-111111111111"),
		Formula:     "xd6",
	})
	require.ErrorIs(t, err, diceRollerServices.ErrInvalidExpression)
}
