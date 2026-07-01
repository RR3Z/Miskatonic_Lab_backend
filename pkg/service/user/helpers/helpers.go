package userHelpers

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
)

func ResolveUsername(input model.UpsertUserInput) string {
	if input.Username != nil && strings.TrimSpace(*input.Username) != "" {
		return strings.TrimSpace(*input.Username)
	}

	email := resolveProvidedEmail(input)
	if email != "" {
		return strings.Split(email, "@")[0]
	}

	hash := sha256.Sum256([]byte(strings.TrimSpace(input.ID)))
	return "user_" + hex.EncodeToString(hash[:])[:12]
}

func ResolveEmail(input model.UpsertUserInput) string {
	email := resolveProvidedEmail(input)
	if email != "" {
		return email
	}

	return input.ID + "@users.local"
}

func resolveProvidedEmail(input model.UpsertUserInput) string {
	if input.PrimaryEmailAddressID != nil {
		for _, e := range input.EmailAddresses {
			if e.ID == *input.PrimaryEmailAddressID {
				return e.EmailAddress
			}
		}
	}

	if len(input.EmailAddresses) > 0 {
		return input.EmailAddresses[0].EmailAddress
	}

	return ""
}
