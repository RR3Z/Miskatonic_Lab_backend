package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func newRoomHandlerTestRouter(roomService room.IRoom) http.Handler {
	h := handler.NewHandler(&service.Service{Room: roomService})
	return h.InitRoutesWithAuth(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{Subject: "user_1"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func performRoomRequest(router http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func testRoomUnitUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
