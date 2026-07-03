package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

const (
	diceTestUserID      = "user_dice_test"
	diceTestCharacterID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	diceTestRoomID      = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
)

func newEventPublishingTestSubject() (*FakeDiceRollerService, *FakeEventPublisher, *diceRollerServices.EventPublishingDiceRollerService) {
	next := &FakeDiceRollerService{Roll: testDiceRollModel(), Rolls: []diceRollerDTO.DiceRollModel{testDiceRollModel()}}
	publisher := &FakeEventPublisher{}
	return next, publisher, diceRollerServices.NewEventPublishingDiceRollerService(next, publisher)
}

func testDiceRollModel() diceRollerDTO.DiceRollModel {
	return diceRollerDTO.DiceRollModel{
		ID:          diceTestUUID(diceTestCharacterID),
		CharacterID: diceTestUUID(diceTestCharacterID),
		UserID:      diceTestUserID,
		Expression:  "2d6+3",
		Result:      10,
		Details:     []byte(`[{"type":"dice","sides":6,"rolls":[3,4]},{"type":"modifier","value":3}]`),
	}
}

func diceTestUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}
	return uuid
}

func requirePublishedEvent(t *testing.T, publisher *FakeEventPublisher, expected events.Event) {
	t.Helper()
	require.Len(t, publisher.Events, 1)
	require.Equal(t, expected, publisher.Events[0])
}

func diceTestMakeRollInput() diceRollerDTO.MakeRollInput {
	return diceRollerDTO.MakeRollInput{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
		Formula:     "2d6+3",
	}
}

func diceTestMakeRollInputWithRoomID() diceRollerDTO.MakeRollInput {
	roomID := diceTestUUID(diceTestRoomID)
	input := diceTestMakeRollInput()
	input.RoomID = &roomID
	return input
}

func diceTestGetLastDiceRollsInput() diceRollerDTO.GetLastDiceRollsInput {
	return diceRollerDTO.GetLastDiceRollsInput{
		UserID:      diceTestUserID,
		CharacterID: diceTestUUID(diceTestCharacterID),
	}
}
