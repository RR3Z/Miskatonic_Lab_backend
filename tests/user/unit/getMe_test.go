package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

func newAuthenticatedMeTestSubject(t *testing.T, userService *FakeUserService) http.Handler {
	t.Helper()

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{
					Subject: "user_test_123",
				},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return handler.NewHandler(&service.Service{User: userService}).InitRoutesWithAuth(authMiddleware)
}

func TestGetMeReturnsUserModel(t *testing.T) {
	userService := &FakeUserService{
		GetUserResult: model.UserModel{
			ID:       "user_test_123",
			Username: "testuser",
			Email:    "test@example.com",
		},
	}
	router := newAuthenticatedMeTestSubject(t, userService)

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, userService.GetUserCalls)
	require.Equal(t, "user_test_123", userService.LastGetUserInput.ID)
	require.JSONEq(t, `{
		"id":"user_test_123",
		"username":"testuser",
		"email":"test@example.com",
		"avatar_url":null,
		"created_at":null,
		"updated_at":null
	}`, recorder.Body.String())
}

func TestGetMeReturnsNotFoundWhenUserMissing(t *testing.T) {
	userService := &FakeUserService{
		GetUserErr: userErrors.ErrUserNotFound,
	}
	router := newAuthenticatedMeTestSubject(t, userService)

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.JSONEq(t, `{"code":"user.not_found","message":"user not found"}`, recorder.Body.String())
}

func TestGetMeReturnsInternalErrorWhenServiceFails(t *testing.T) {
	userService := &FakeUserService{}
	router := newAuthenticatedMeTestSubject(t, userService)
	userService.GetUserErr = http.ErrServerClosed

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.JSONEq(t, `{"code":"common.internal_error","message":"failed to get user"}`, recorder.Body.String())
}
