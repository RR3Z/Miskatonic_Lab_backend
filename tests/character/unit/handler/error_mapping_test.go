package tests

import (
	"errors"
	"net/http"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

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
			_, router := newCharacterHandlerTestSubject(tc.err)

			recorder := performCharacterRequest(router, tc.method, tc.path, tc.body)

			requireCharacterError(t, recorder, tc.wantStatus, tc.wantCode)
		})
	}
}
