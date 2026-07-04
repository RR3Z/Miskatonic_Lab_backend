package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2EMultiUserRoomSelectedCharacterVisibility(t *testing.T) {
	owner := newE2ESubject(t)
	player := newSecondE2ESubject(t)
	if owner.userID == player.userID {
		t.Skip("E2E_SECOND_AUTH_TOKEN must belong to a different Clerk subject")
	}

	password := "e2e-multi-" + e2eHash(owner.userID+player.userID)
	room := owner.createRoom(t, password)
	require.Equal(t, owner.userID, room.OwnerID)
	t.Cleanup(func() {
		owner.deleteRoom(t, room.ID)
	})

	ownerCharacterID := owner.createCharacter(t, "E2E Keeper Investigator")
	t.Cleanup(func() {
		owner.deleteCharacter(t, ownerCharacterID)
	})
	playerCharacterID := player.createCharacter(t, "E2E Player Investigator")
	t.Cleanup(func() {
		player.deleteCharacter(t, playerCharacterID)
	})

	player.joinRoom(t, room.ID, password)
	owner.selectCharacter(t, room.ID, ownerCharacterID)
	player.selectCharacter(t, room.ID, playerCharacterID)

	var ownerVisible []e2eSelectedCharacterResponse
	owner.doJSON(t, http.MethodGet, "/api/rooms/"+url.PathEscape(room.ID)+"/characters", nil, http.StatusOK, &ownerVisible)
	require.ElementsMatch(t, []string{ownerCharacterID, playerCharacterID}, collectSelectedCharacterIDs(ownerVisible))

	var playerVisible []e2eSelectedCharacterResponse
	player.doJSON(t, http.MethodGet, "/api/rooms/"+url.PathEscape(room.ID)+"/characters", nil, http.StatusOK, &playerVisible)
	require.Equal(t, []string{playerCharacterID}, collectSelectedCharacterIDs(playerVisible))
}
