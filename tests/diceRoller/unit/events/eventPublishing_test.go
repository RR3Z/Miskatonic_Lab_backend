package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/stretchr/testify/require"
)

var errDiceTest = errors.New("dice roll failed")

func TestEventPublishingDiceRoller_MakeRoll_Success(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()

	roll, err := svc.MakeRoll(context.Background(), diceRollerServices.DiceRollInput{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
		Formula:     "2d6+3",
	})
	require.NoError(t, err)
	require.Equal(t, next.Roll, roll)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Expression:  "2d6+3",
		Result:      10,
	})
}

func TestEventPublishingDiceRoller_MakeRoll_Failure(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Err = errDiceTest

	_, err := svc.MakeRoll(context.Background(), diceRollerServices.DiceRollInput{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
		Formula:     "2d6+3",
	})
	require.ErrorIs(t, err, errDiceTest)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeFailed{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Err:         errDiceTest,
	})
}

func TestEventPublishingDiceRoller_GetLastDiceRolls_Success(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()

	rolls, err := svc.GetLastDiceRolls(context.Background(), dbGetDiceRollsParams())
	require.NoError(t, err)
	require.Equal(t, next.Rolls, rolls)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollsListSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Count:       1,
	})
}

func TestEventPublishingDiceRoller_GetLastDiceRolls_Empty(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Rolls = []db.DiceRoll{}

	rolls, err := svc.GetLastDiceRolls(context.Background(), dbGetDiceRollsParams())
	require.NoError(t, err)
	require.Empty(t, rolls)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollsListSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Count:       0,
	})
}

func TestEventPublishingDiceRoller_GetLastDiceRolls_Failure(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Err = errDiceTest

	_, err := svc.GetLastDiceRolls(context.Background(), dbGetDiceRollsParams())
	require.ErrorIs(t, err, errDiceTest)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollsListFailed{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Err:         errDiceTest,
	})
}

func TestEventPublishingDiceRoller_MultipleCallsAccumulateEvents(t *testing.T) {
	_, publisher, svc := newEventPublishingTestSubject()

	_, err := svc.MakeRoll(context.Background(), diceRollerServices.DiceRollInput{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
		Formula:     "1d20",
	})
	require.NoError(t, err)

	_, err = svc.GetLastDiceRolls(context.Background(), dbGetDiceRollsParams())
	require.NoError(t, err)

	require.Len(t, publisher.Events, 2)
	require.IsType(t, diceEvents.DiceRollMakeSucceeded{}, publisher.Events[0])
	require.IsType(t, diceEvents.DiceRollsListSucceeded{}, publisher.Events[1])
}

func dbGetDiceRollsParams() db.GetDiceRollsParams {
	return db.GetDiceRollsParams{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
	}
}
