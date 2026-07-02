package roomHelpers

import "strings"

func OptionalPassword(password string) *string {
	if strings.TrimSpace(password) == "" {
		return nil
	}
	return &password
}
