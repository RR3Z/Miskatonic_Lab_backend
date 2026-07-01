package user

import (
	"strings"

	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
)

func validateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return userErrors.ErrMissingUserID
	}
	return nil
}
