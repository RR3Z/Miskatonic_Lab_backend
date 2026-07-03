package roomHelpers

import (
	"encoding/json"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
)

func ChatMessagePayload(text string) ([]byte, error) {
	return json.Marshal(roomEvents.ChatMessagePayload{
		Text: text,
	})
}

func OwnerTransferredPayload(previousOwnerID string, newOwnerID string) ([]byte, error) {
	return json.Marshal(roomEvents.OwnerTransferredPayload{
		PreviousOwnerID: previousOwnerID,
		NewOwnerID:      newOwnerID,
	})
}

func DiceRollPayload(rollID, characterID, expression string, result int32, details []byte) ([]byte, error) {
	return json.Marshal(roomEvents.DiceRollPayload{
		RollID:      rollID,
		CharacterID: characterID,
		Expression:  expression,
		Result:      result,
		Details:     details,
	})
}
