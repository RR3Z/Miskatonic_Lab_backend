package roomHelpers

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func passwordMatches(hash string, password string) bool {
	if strings.TrimSpace(hash) == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(strings.TrimSpace(password)))
	return err == nil
}
