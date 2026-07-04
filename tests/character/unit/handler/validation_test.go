package tests

import (
	"net/http"
	"testing"
)

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
			requireRejectedBeforeService(t, http.MethodGet, tc.path, "", "character.invalid_id")
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
			requireRejectedBeforeService(t, tc.method, tc.path, `{"broken"`, "character.invalid_input")
		})
	}
}
