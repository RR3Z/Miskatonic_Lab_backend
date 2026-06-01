package middleware

import (
	"net/http"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(next)
}
