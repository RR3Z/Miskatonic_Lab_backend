package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newRoomHandlerTestRouterWithHub(roomService room.IRoom, hub *ws.RoomHub) http.Handler {
	handler := roomHandler.NewWithHub(roomService, hub)
	router := chi.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{Subject: "user_1"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	router.Route("/api/rooms", handler.RegisterRoutes)
	return router
}

func registerRoomUnitTestClient(t *testing.T, hub *ws.RoomHub, roomID pgtype.UUID) <-chan roomEvents.Event {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go hub.Run(ctx)

	client, events := ws.NewTestClientWithUser(hub, roomID.String(), "user_1")
	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	return events
}

func requireRoomUnitClientClosed(t *testing.T, events <-chan roomEvents.Event) {
	t.Helper()

	require.Eventually(t, func() bool {
		select {
		case _, ok := <-events:
			return !ok
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)
}
