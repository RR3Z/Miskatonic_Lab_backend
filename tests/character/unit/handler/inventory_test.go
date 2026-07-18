package tests

import (
	"net/http"
	"testing"
)

func TestInventoryItemRejectsInvalidItemIDBeforeService(t *testing.T) {
	requireRejectedBeforeService(t, http.MethodGet, "/api/characters/"+testCharacterID+"/inventory/not-a-uuid/", "", "character.invalid_id")
}

func TestInventoryItemRejectsInvalidCharacterIDBeforeService(t *testing.T) {
	requireRejectedBeforeService(t, http.MethodPost, "/api/characters/not-a-uuid/inventory/", `{ "name": "Flashlight" }`, "character.invalid_id")
}
