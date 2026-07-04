package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	testCharacterID = "11111111-1111-1111-1111-111111111111"
	testItemID      = "22222222-2222-2222-2222-222222222222"
	testSkillID     = "33333333-3333-3333-3333-333333333333"
	testNoteID      = "44444444-4444-4444-4444-444444444444"
)

func newCharacterHandlerTestRouter(characterService *fakeCharacterHandlerService) http.Handler {
	h := handler.NewHandler(&service.Service{Character: characterService})
	return h.InitRoutesWithAuth(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{Subject: "user_1"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func performCharacterRequest(router http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func testCharacterUnitUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
