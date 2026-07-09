package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
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

func newCharacterHandlerTestSubject(err error) (*fakeCharacterHandlerService, http.Handler) {
	service := &fakeCharacterHandlerService{err: err}
	return service, newCharacterHandlerTestRouter(service)
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

func requireCharacterError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
	t.Helper()

	require.Equal(t, status, recorder.Code)
	require.Contains(t, recorder.Body.String(), code)
}

func requireRejectedBeforeService(t *testing.T, method string, path string, body string, code string) {
	t.Helper()

	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(router, method, path, body)

	requireCharacterError(t, recorder, http.StatusBadRequest, code)
	require.Zero(t, service.totalCalls())
}

func testCharacterUnitUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}

func characterString(value string) *string {
	return &value
}

func characterInt16(value int16) *int16 {
	return &value
}
