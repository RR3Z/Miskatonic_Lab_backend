package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

func TestGetUserIDFromContextReturnsClerkSubject(t *testing.T) {
	ctx := clerk.ContextWithSessionClaims(context.Background(), &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{Subject: "user_1"},
	})

	userID := utils.GetUserIDFromContext(ctx)

	require.Equal(t, "user_1", userID)
}

func TestGetUserIDFromContextReturnsEmptyStringForEmptySubject(t *testing.T) {
	ctx := clerk.ContextWithSessionClaims(context.Background(), &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{Subject: ""},
	})

	userID := utils.GetUserIDFromContext(ctx)

	require.Empty(t, userID)
}

func TestGetUserIDFromContextPanicsWhenClaimsAreMissing(t *testing.T) {
	require.Panics(t, func() {
		_ = utils.GetUserIDFromContext(context.Background())
	})
}
