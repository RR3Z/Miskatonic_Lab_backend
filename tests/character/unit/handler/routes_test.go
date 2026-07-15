package tests

import (
	"net/http"
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/stretchr/testify/require"
)

func TestCharacterRoutesAreMounted(t *testing.T) {
	characterID := testCharacterUnitUUID(testCharacterID)
	service := &fakeCharacterHandlerService{
		characters: []characterDTO.CharacterSummaryModel{{ID: characterID, Name: "Investigator"}},
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
		{"get derived stats", http.MethodGet, "/api/characters/" + testCharacterID + "/derived-stats/"},
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

func TestCharacterDerivedStatsRoutesAreReadOnly(t *testing.T) {
	router := newCharacterHandlerTestRouter(&fakeCharacterHandlerService{})
	path := "/api/characters/" + testCharacterID + "/derived-stats/"

	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			recorder := performCharacterRequest(router, method, path, `{}`)

			require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
		})
	}
}
