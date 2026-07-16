package tests

import (
	"encoding/json"
	"testing"

	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
	"github.com/stretchr/testify/require"
)

func TestDiceRollPayloadIncludesStructuredDetails(t *testing.T) {
	payloadJSON, err := roomHelpers.DiceRollPayload(
		"roll-1",
		"character-1",
		"1d100",
		24,
		[]byte(`{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`),
	)
	require.NoError(t, err)

	var payload roomModel.DiceRollPayload
	require.NoError(t, json.Unmarshal(payloadJSON, &payload))
	require.JSONEq(t, `{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`, string(payload.Details))
}
