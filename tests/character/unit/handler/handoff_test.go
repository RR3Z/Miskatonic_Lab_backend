package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/stretchr/testify/require"
)

func TestCreateCharacterPassesBodyAndUserID(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/", `{
		"name":"Harvey Walters",
		"occupation":"Professor",
		"age":42,
		"sex":"male",
		"residence":"Arkham",
		"birthplace":"Boston"
	}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, service.createCalls)
	require.Equal(t, "user_1", service.createInput.UserID)
	require.Equal(t, "Harvey Walters", service.createInput.Name)
	require.Equal(t, "Professor", *service.createInput.Occupation)
	require.Equal(t, int16(42), *service.createInput.Age)
	require.Equal(t, "Boston", *service.createInput.Birthplace)
}

func TestCreateCharacterRejectsRemovedPlayerName(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/", `{
		"name":"Harvey Walters",
		"player_name":"Roger"
	}`)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.invalid_input")
	require.Zero(t, service.createCalls)
}

func TestCreateCharacterRejectsClientPortraitURL(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/", `{
		"name":"Harvey Walters",
		"portrait_url":"https://assets.example.test/portraits/harvey.webp"
	}`)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.portrait_managed_by_server")
	require.Zero(t, service.createCalls)
}

func TestCreateCharacterRejectsNullClientPortraitURL(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPost, "/api/characters/", `{
		"name":"Harvey Walters",
		"portrait_url":null
	}`)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.portrait_managed_by_server")
	require.Zero(t, service.createCalls)
}

func TestUpdateCharacterRejectsClientPortraitURL(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, http.MethodPut, "/api/characters/"+testCharacterID+"/", `{
		"name":"Harvey Walters",
		"portrait_url":"https://assets.example.test/portraits/updated.webp"
	}`)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.portrait_managed_by_server")
	require.Zero(t, service.updateCalls)
}

func TestGetAllCharactersReturnsSummaryJSON(t *testing.T) {
	portraitURL := "https://assets.example.test/portraits/summary.webp"
	card := characterDTO.CharacterSummaryModel{
		ID:          testCharacterUnitUUID(testCharacterID),
		Name:        "Harvey Walters",
		Occupation:  characterString("Professor"),
		Age:         characterInt16(42),
		Sex:         characterString(""),
		Residence:   characterString("Arkham"),
		PortraitUrl: &portraitURL,
	}
	card.HP.Current = 7
	card.HP.Max = 12
	card.MP.Current = 4
	card.MP.Max = 9
	card.Sanity.Current = 33
	card.Sanity.Max = 60
	card.Luck.Current = 20
	card.Luck.Starting = 45

	service := &fakeCharacterHandlerService{characters: []characterDTO.CharacterSummaryModel{card}}
	router := newCharacterHandlerTestRouter(service)

	recorder := performCharacterRequest(router, http.MethodGet, "/api/characters/", "")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, service.getAllCalls)
	require.Equal(t, "user_1", service.getAllUserID)

	var response []map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Len(t, response, 1)
	require.Equal(t, "Harvey Walters", response[0]["name"])
	require.Equal(t, "Professor", response[0]["occupation"])
	require.Equal(t, "", response[0]["sex"])
	require.Equal(t, "Arkham", response[0]["residence"])
	require.Equal(t, portraitURL, response[0]["portrait_url"])
	require.NotContains(t, response[0], "birthplace")
	require.NotContains(t, response[0], "user_id")
	require.NotContains(t, response[0], "created_at")
	require.NotContains(t, response[0], "updated_at")

	hp := response[0]["hp"].(map[string]any)
	require.Equal(t, float64(7), hp["current_hp"])
	require.Equal(t, float64(12), hp["max_hp"])
	mp := response[0]["mp"].(map[string]any)
	require.Equal(t, float64(4), mp["current_mp"])
	require.Equal(t, float64(9), mp["max_mp"])
	sanity := response[0]["sanity"].(map[string]any)
	require.Equal(t, float64(33), sanity["current_sanity"])
	require.Equal(t, float64(60), sanity["max_sanity"])
	luck := response[0]["luck"].(map[string]any)
	require.Equal(t, float64(20), luck["current_luck"])
	require.Equal(t, float64(45), luck["starting_luck"])
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
