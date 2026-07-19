package tests

import (
	"encoding/json"
	"testing"

	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestRoomModelExposesInviteTokenOnlyToOwner(t *testing.T) {
	room := db.Room{
		InviteToken: "invite-secret",
		OwnerID:     "owner-1",
	}

	ownerRoom := roomModels.ToRoomModel(room, nil, "owner-1")
	memberRoom := roomModels.ToRoomModel(room, nil, "member-1")

	require.Equal(t, "invite-secret", ownerRoom.InviteToken)
	require.Empty(t, memberRoom.InviteToken)

	memberJSON, err := json.Marshal(memberRoom)
	require.NoError(t, err)
	require.NotContains(t, string(memberJSON), "invite_token")
}
