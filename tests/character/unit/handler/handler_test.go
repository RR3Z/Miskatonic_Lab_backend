package tests

import (
	"errors"
	"net/http"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCharacterRoutesAreMounted(t *testing.T) {
	characterID := testCharacterUnitUUID(testCharacterID)
	service := &fakeCharacterHandlerService{
		characters: []characterDTO.CharacterShortModel{{ID: characterID, UserID: "user_1", Name: "Investigator"}},
		character:  characterDTO.CharacterModel{CharacterShortModel: characterDTO.CharacterShortModel{ID: characterID, UserID: "user_1", Name: "Investigator"}},
	}
	router := newCharacterHandlerTestRouter(service)

	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"list characters", http.MethodGet, "/api/characters/"},
		{"get character", http.MethodGet, "/api/characters/" + testCharacterID + "/"},
		{"get health", http.MethodGet, "/api/characters/" + testCharacterID + "/health/"},
		{"get characteristics", http.MethodGet, "/api/characters/" + testCharacterID + "/characteristics/"},
		{"get backstory item", http.MethodGet, "/api/characters/" + testCharacterID + "/backstory/items/" + testItemID + "/"},
		{"get skill", http.MethodGet, "/api/characters/" + testCharacterID + "/skills/" + testSkillID + "/"},
		{"get note", http.MethodGet, "/api/characters/" + testCharacterID + "/notes/" + testNoteID + "/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCharacterRequest(router, tc.method, tc.path, "")

			require.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}

func TestCreateCharacterPassesBodyAndUserID(t *testing.T) {
	service := &fakeCharacterHandlerService{}
	router := newCharacterHandlerTestRouter(service)

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
	service := &fakeCharacterHandlerService{}
	router := newCharacterHandlerTestRouter(service)

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
	service := &fakeCharacterHandlerService{}
	router := newCharacterHandlerTestRouter(service)

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
	service := &fakeCharacterHandlerService{}
	router := newCharacterHandlerTestRouter(service)

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
	service := &fakeCharacterHandlerService{}
	router := newCharacterHandlerTestRouter(service)

	recorder := performCharacterRequest(router, http.MethodDelete, "/api/characters/"+testCharacterID+"/", "")

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Empty(t, recorder.Body.String())
	require.Equal(t, 1, service.deleteCalls)
	require.Equal(t, "user_1", service.deleteInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.deleteInput.ID)
}

func TestCharacterRoutesRejectInvalidUUIDBeforeService(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{"character id", "/api/characters/not-a-uuid/"},
		{"backstory item id", "/api/characters/" + testCharacterID + "/backstory/items/not-a-uuid/"},
		{"skill id", "/api/characters/" + testCharacterID + "/skills/not-a-uuid/"},
		{"note id", "/api/characters/" + testCharacterID + "/notes/not-a-uuid/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &fakeCharacterHandlerService{}
			router := newCharacterHandlerTestRouter(service)

			recorder := performCharacterRequest(router, http.MethodGet, tc.path, "")

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Contains(t, recorder.Body.String(), "character.invalid_id")
			require.Zero(t, service.totalCalls())
		})
	}
}

func TestCharacterRoutesRejectInvalidJSONBeforeService(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"create character", http.MethodPost, "/api/characters/"},
		{"upsert health", http.MethodPut, "/api/characters/" + testCharacterID + "/health/"},
		{"create note", http.MethodPost, "/api/characters/" + testCharacterID + "/notes/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &fakeCharacterHandlerService{}
			router := newCharacterHandlerTestRouter(service)

			recorder := performCharacterRequest(router, tc.method, tc.path, `{"broken"`)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Contains(t, recorder.Body.String(), "character.invalid_input")
			require.Zero(t, service.totalCalls())
		})
	}
}

func TestCharacterServiceErrorsMapToJSON(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "not found",
			method:     http.MethodGet,
			path:       "/api/characters/" + testCharacterID + "/",
			err:        pgx.ErrNoRows,
			wantStatus: http.StatusNotFound,
			wantCode:   "character.not_found",
		},
		{
			name:       "name required",
			method:     http.MethodPost,
			path:       "/api/characters/",
			body:       `{"name":""}`,
			err:        characterErrors.ErrNameRequired,
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.name_required",
		},
		{
			name:       "state current exceeds max",
			method:     http.MethodPut,
			path:       "/api/characters/" + testCharacterID + "/health/",
			body:       `{"max_hp":1,"current_hp":2}`,
			err:        myErrors.ErrCurrentHealthExceedsMax,
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.state_current_exceeds_max",
		},
		{
			name:       "skill in use",
			method:     http.MethodDelete,
			path:       "/api/characters/" + testCharacterID + "/skills/" + testSkillID + "/",
			err:        characterErrors.ErrSkillInUse,
			wantStatus: http.StatusConflict,
			wantCode:   "character.skill_in_use",
		},
		{
			name:       "fallback internal error",
			method:     http.MethodGet,
			path:       "/api/characters/",
			err:        errors.New("repository unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "common.internal_error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &fakeCharacterHandlerService{err: tc.err}
			router := newCharacterHandlerTestRouter(service)

			recorder := performCharacterRequest(router, tc.method, tc.path, tc.body)

			require.Equal(t, tc.wantStatus, recorder.Code)
			require.Contains(t, recorder.Body.String(), tc.wantCode)
		})
	}
}
