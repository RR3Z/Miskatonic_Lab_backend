package room

import "strings"

func validateMaxPlayers(maxPlayers int32) error {
	if maxPlayers < 1 {
		return ErrInvalidInput
	}
	return nil
}

func validateInviteToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidInput
	}
	return nil
}

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrInvalidPassword
	}
	return nil
}
