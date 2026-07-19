package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestRoomRoutesRequireAuthentication(t *testing.T) {
	router := handler.NewHandler(handler.Dependencies{
		Services: &service.Service{},
	}).InitRoutes(middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{}))

	roomID := "11111111-1111-1111-1111-111111111111"
	requests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "list", method: http.MethodGet, path: "/api/rooms/"},
		{name: "create", method: http.MethodPost, path: "/api/rooms/"},
		{name: "get", method: http.MethodGet, path: "/api/rooms/" + roomID + "/"},
		{name: "update", method: http.MethodPut, path: "/api/rooms/" + roomID + "/"},
		{name: "delete", method: http.MethodDelete, path: "/api/rooms/" + roomID + "/"},
		{name: "selected characters", method: http.MethodGet, path: "/api/rooms/" + roomID + "/characters"},
		{name: "events", method: http.MethodGet, path: "/api/rooms/" + roomID + "/events"},
		{name: "websocket", method: http.MethodGet, path: "/api/rooms/" + roomID + "/ws"},
		{name: "transfer ownership", method: http.MethodPut, path: "/api/rooms/" + roomID + "/owner"},
		{name: "join", method: http.MethodPost, path: "/api/rooms/" + roomID + "/join"},
		{name: "leave", method: http.MethodDelete, path: "/api/rooms/" + roomID + "/leave"},
		{name: "kick", method: http.MethodDelete, path: "/api/rooms/" + roomID + "/kick/user-2"},
		{name: "select character", method: http.MethodPut, path: "/api/rooms/" + roomID + "/character"},
		{name: "change role", method: http.MethodPut, path: "/api/rooms/" + roomID + "/members/user-2/role"},
	}

	for _, test := range requests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader("{}"))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusForbidden, recorder.Code)
		})
	}
}
