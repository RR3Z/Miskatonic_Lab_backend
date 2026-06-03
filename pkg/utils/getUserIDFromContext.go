package utils

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		slog.Error(
			"failed to get clerk session claims",
			"component", "character_api",
		)
		return "", errors.New("unauthorized")
	}

	userID := claims.Subject
	if strings.TrimSpace(userID) == "" {
		slog.Error(
			"clerk session claims missing subject",
			"component", "user_api",
		)
		return "", errors.New("unauthorized")
	}

	return userID, nil
}
