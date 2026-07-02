package roomHelpers

import (
	"strings"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

func HasAnyJoinCredential(input model.JoinRoomInput) bool {
	return hasJoinCredential(input.InviteToken) || hasJoinCredential(input.Password)
}

func CanUseJoinCredential(inviteToken string, passwordHash string, input model.JoinRoomInput) bool {
	if hasJoinCredential(input.InviteToken) && input.InviteToken == inviteToken {
		return true
	}

	if hasJoinCredential(input.Password) && passwordMatches(passwordHash, input.Password) {
		return true
	}

	return false
}

func hasJoinCredential(input string) bool {
	return strings.TrimSpace(input) != ""
}
