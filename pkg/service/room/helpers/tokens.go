package roomHelpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateInviteToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
