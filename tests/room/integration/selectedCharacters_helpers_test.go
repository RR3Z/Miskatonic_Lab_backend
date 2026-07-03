package tests

import (
	"testing"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/stretchr/testify/require"
)

func requireSelectedCharacterUsers(t *testing.T, characters []model.SelectedCharacterModel, expectedUsers ...string) {
	t.Helper()

	users := selectedCharacterUsers(characters)
	for _, userID := range expectedUsers {
		require.Contains(t, users, userID)
	}
}

func selectedCharacterUsers(characters []model.SelectedCharacterModel) []string {
	users := make([]string, 0, len(characters))
	for _, character := range characters {
		users = append(users, character.UserID)
	}
	return users
}
