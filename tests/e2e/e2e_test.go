package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestE2EProtectedRouteRejectsMissingToken(t *testing.T) {
	requireE2EEnabled(t)

	baseURL := e2eBaseURL()
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/me", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, res.StatusCode)
}

func TestE2EGetMeReturnsCurrentUser(t *testing.T) {
	subject := newE2ESubject(t)

	var user e2eUserResponse
	subject.doJSON(t, http.MethodGet, "/api/me", nil, http.StatusOK, &user)

	require.Equal(t, subject.userID, user.ID)
	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.Email)
}

func TestE2ECharacterHTTPFlow(t *testing.T) {
	subject := newE2ESubject(t)

	characterID := subject.createCharacter(t, "E2E Investigator")
	t.Cleanup(func() {
		subject.deleteCharacter(t, characterID)
	})

	var characters []e2eIDResponse
	subject.doJSON(t, http.MethodGet, "/api/characters/", nil, http.StatusOK, &characters)
	require.Contains(t, collectIDs(characters), characterID)

	var health e2eHealthResponse
	subject.doJSON(
		t,
		http.MethodPut,
		"/api/characters/"+characterID+"/health/",
		map[string]any{"max_hp": 12, "current_hp": 7},
		http.StatusOK,
		&health,
	)
	require.Equal(t, int16(12), health.MaxHp)
	require.Equal(t, int16(7), health.CurrentHp)

	var full e2eCharacterResponse
	subject.doJSON(t, http.MethodGet, "/api/characters/"+characterID+"/", nil, http.StatusOK, &full)
	require.Equal(t, characterID, full.ID)
	require.Equal(t, subject.userID, full.UserID)
	require.Equal(t, "E2E Investigator", full.Name)
	require.Equal(t, int16(12), full.HP.MaxHp)
	require.Equal(t, int16(7), full.HP.CurrentHp)

	subject.doJSON(t, http.MethodDelete, "/api/characters/"+characterID+"/", nil, http.StatusNoContent, nil)
}

func TestE2ECharacterOwnershipDenial(t *testing.T) {
	subject := newE2ESubject(t)

	otherUserID := "e2e_other_" + e2eHash(time.Now().String())
	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:       otherUserID,
		Username: "e2e_other_" + e2eHash(otherUserID),
		Email:    "e2e_other+" + e2eHash(otherUserID) + "@example.com",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = subject.queries.DeleteUserByClerkID(context.Background(), otherUserID)
	})

	otherCharacter, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{
		UserID: otherUserID,
		Name:   "Other Investigator",
	})
	require.NoError(t, err)

	subject.doJSON(t, http.MethodGet, "/api/characters/"+otherCharacter.ID.String()+"/", nil, http.StatusNotFound, nil)
}

func TestE2ERoomDiceRollCreatesRoomEvent(t *testing.T) {
	subject := newE2ESubject(t)

	characterID := subject.createCharacter(t, "E2E Room Investigator")
	t.Cleanup(func() {
		subject.deleteCharacter(t, characterID)
	})

	room := subject.createRoom(t, "e2e-room-password")
	require.Equal(t, subject.userID, room.OwnerID)
	t.Cleanup(func() {
		subject.deleteRoom(t, room.ID)
	})

	subject.selectCharacter(t, room.ID, characterID)

	var roll e2eIDResponse
	subject.doJSON(
		t,
		http.MethodPost,
		"/api/dice-roll/"+characterID+"/",
		map[string]any{"expression": "1d6", "room_id": room.ID},
		http.StatusCreated,
		&roll,
	)
	require.NotEmpty(t, roll.ID)

	events := subject.waitForRoomEvents(t, room.ID, "dice.roll")
	require.NotEmpty(t, events)
}
