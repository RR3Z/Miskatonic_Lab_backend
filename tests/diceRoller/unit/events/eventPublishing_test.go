package tests

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/stretchr/testify/require"
)

var errDiceTest = errors.New("dice roll failed")

func TestEventPublishingDiceRoller_MakeRoll_Success(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()

	roll, err := svc.MakeRoll(context.Background(), diceTestMakeRollInput())
	require.NoError(t, err)
	require.Equal(t, next.Roll, roll)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		RollID:      diceTestCharacterID,
		Expression:  "2d6+3",
		Result:      10,
		Details:     []byte(`{"rolls":[{"type":"dice","sides":6,"rolls":[3,4]},{"type":"modifier","value":3}]}`),
	})
}

func TestEventPublishingDiceRoller_MakeRoll_SuccessWithRoomID(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	roomID := diceTestRoomID

	roll, err := svc.MakeRoll(context.Background(), diceTestMakeRollInputWithRoomID())
	require.NoError(t, err)
	require.Equal(t, next.Roll, roll)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		RollID:      diceTestCharacterID,
		Expression:  "2d6+3",
		Result:      10,
		Details:     []byte(`{"rolls":[{"type":"dice","sides":6,"rolls":[3,4]},{"type":"modifier","value":3}]}`),
		RoomID:      &roomID,
	})
}

func TestEventPublishingDiceRoller_MakeRoll_ForwardsStructuredDetails(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Roll = diceRollerDTO.DiceRollModel{
		ID:          diceTestUUID(diceTestCharacterID),
		CharacterID: diceTestUUID(diceTestCharacterID),
		UserID:      diceTestUserID,
		Expression:  "1d100",
		Result:      24,
		Details:     json.RawMessage(`{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`),
	}

	_, err := svc.MakeRoll(context.Background(), diceTestMakeRollInput())
	require.NoError(t, err)

	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeSucceeded{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		RollID:      diceTestCharacterID,
		Expression:  "1d100",
		Result:      24,
		Details:     next.Roll.Details,
	})
}

func TestEventPublishingDiceRoller_MakeRoll_Failure(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Err = errDiceTest

	_, err := svc.MakeRoll(context.Background(), diceTestMakeRollInput())
	require.ErrorIs(t, err, errDiceTest)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeFailed{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Err:         errDiceTest,
	})
}

func TestEventPublishingDiceRoller_MakeRoll_FailureWithRoomID(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()
	next.Err = errDiceTest
	roomID := diceTestRoomID

	_, err := svc.MakeRoll(context.Background(), diceTestMakeRollInputWithRoomID())
	require.ErrorIs(t, err, errDiceTest)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollMakeFailed{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Err:         errDiceTest,
		RoomID:      &roomID,
	})
}

func TestEventPublishingDiceRoller_GetLastDiceRolls_Success(t *testing.T) {
	next, publisher, svc := newEventPublishingTestSubject()

	rolls, err := svc.GetLastDiceRolls(context.Background(), diceTestGetLastDiceRollsInput())
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
	next.Rolls = nil

	rolls, err := svc.GetLastDiceRolls(context.Background(), diceTestGetLastDiceRollsInput())
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

	_, err := svc.GetLastDiceRolls(context.Background(), diceTestGetLastDiceRollsInput())
	require.ErrorIs(t, err, errDiceTest)
	requirePublishedEvent(t, publisher, diceEvents.DiceRollsListFailed{
		UserID:      diceTestUserID,
		CharacterID: diceTestCharacterID,
		Err:         errDiceTest,
	})
}

func TestEventPublishingDiceRoller_MultipleCallsAccumulateEvents(t *testing.T) {
	_, publisher, svc := newEventPublishingTestSubject()

	_, err := svc.MakeRoll(context.Background(), diceTestMakeRollInput())
	require.NoError(t, err)

	_, err = svc.GetLastDiceRolls(context.Background(), diceTestGetLastDiceRollsInput())
	require.NoError(t, err)

	require.Len(t, publisher.Events, 2)
	require.IsType(t, diceEvents.DiceRollMakeSucceeded{}, publisher.Events[0])
	require.IsType(t, diceEvents.DiceRollsListSucceeded{}, publisher.Events[1])
}
