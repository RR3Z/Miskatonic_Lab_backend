package utils

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
)

func GetUserIDFromContext(ctx context.Context) string {
	claims, _ := clerk.SessionClaimsFromContext(ctx)
	return claims.Subject
}
