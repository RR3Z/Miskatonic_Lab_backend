package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"

	diceRollerHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-chi/chi/v5"
)

func newDiceRollerTestRouter(svc *fakeDiceRollerHandlerService) http.Handler {
	services := &service.Service{
		DiceRoller: svc,
	}
	h := handler.NewHandler(services)
	return h.InitRoutesWithAuth(diceRollerAuthMiddleware)
}

func newDiceRollerTestRouterWithRoom(svc *fakeDiceRollerHandlerService, checker diceRollerHandler.RoomAccessChecker) http.Handler {
	diceHandler := diceRollerHandler.NewWithRoomChecker(svc, checker)
	router := chi.NewRouter()
	router.Route("/api/dice-roll", diceHandler.RegisterRoutes)
	return diceRollerAuthMiddleware(router)
}

func diceRollerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
			RegisteredClaims: clerk.RegisteredClaims{Subject: "user_test_01"},
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func performDiceRollerRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
