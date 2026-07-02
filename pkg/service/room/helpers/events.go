package roomHelpers

import (
	"encoding/json"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
)

func OwnerTransferredPayload(previousOwnerID string, newOwnerID string) ([]byte, error) {
	return json.Marshal(roomEvents.OwnerTransferredPayload{
		PreviousOwnerID: previousOwnerID,
		NewOwnerID:      newOwnerID,
	})
}
