package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCharacterPassesBodyAndUserID(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/", `{
		"name":"Harvey Walters",
		"player_name":"Roger",
		"occupation":"Professor",
		"age":42,
		"sex":"m",
		"residence":"Arkham",
		"birthplace":"Boston"
	}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, service.createCalls)
	require.Equal(t, "user_1", service.createInput.UserID)
	require.Equal(t, "Harvey Walters", service.createInput.Name)
	require.Equal(t, "Professor", *service.createInput.Occupation)
	require.Equal(t, int16(42), *service.createInput.Age)
}

func TestUpsertHealthPassesPathAndBody(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPut, "/api/characters/"+testCharacterID+"/health/", `{
		"max_hp":12,
		"current_hp":7,
		"major_wound":true,
		"unconscious":false,
		"dying":false,
		"dead":false
	}`)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, service.upsertHealthCalls)
	require.Equal(t, "user_1", service.upsertHealthInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.upsertHealthInput.CharacterID)
	require.Equal(t, int16(12), *service.upsertHealthInput.MaxHp)
	require.Equal(t, int16(7), *service.upsertHealthInput.CurrentHp)
	require.True(t, *service.upsertHealthInput.MajorWound)
}

func TestCreateBackstoryItemPassesPathAndBody(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/"+testCharacterID+"/backstory/items/", `{
		"section":"ideology_beliefs",
		"title":"Old Motto",
		"text":"Trust the archive."
	}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, service.createBackstoryItemCalls)
	require.Equal(t, "user_1", service.createBackstoryItemInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.createBackstoryItemInput.CharacterID)
	require.Equal(t, "ideology_beliefs", service.createBackstoryItemInput.Section)
	require.Equal(t, "Old Motto", service.createBackstoryItemInput.Title)
	require.Equal(t, "Trust the archive.", service.createBackstoryItemInput.Text)
}

func TestCreateNotePassesPathAndBody(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/"+testCharacterID+"/notes/", `{
		"title":"Session One",
		"body":"Found strange tracks."
	}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, service.createNoteCalls)
	require.Equal(t, "user_1", service.createNoteInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.createNoteInput.CharacterID)
	require.Equal(t, "Session One", service.createNoteInput.Title)
	require.Equal(t, "Found strange tracks.", service.createNoteInput.Body)
}

func TestDeleteCharacterReturnsNoContent(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodDelete, "/api/characters/"+testCharacterID+"/", "")

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Empty(t, recorder.Body.String())
	require.Equal(t, 1, service.deleteCalls)
	require.Equal(t, "user_1", service.deleteInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.deleteInput.ID)
}
