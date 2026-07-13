package handler

import (
	"context"
	"net/http"

	characterHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character"
	diceRollerHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller"
	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	userHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/user"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
)

type Dependencies struct {
	Services           *service.Service
	PortraitFileServer http.Handler
}

type Handler struct {
	domainHandlers     *domainHandlers
	portraitFileServer http.Handler
}

type domainHandlers struct {
	characterHandler  *characterHandler.CharacterHandler
	diceRollerHandler *diceRollerHandler.DiceRollerHandler
	roomHandler       *roomHandler.RoomHandler
	userHandler       *userHandler.UserHandler
}

func NewHandler(dependencies Dependencies) *Handler {
	services := dependencies.Services
	roomHub := ws.NewRoomHub()
	go roomHub.Run(context.Background())

	var diceHandler *diceRollerHandler.DiceRollerHandler
	if services.Room != nil {
		diceHandler = diceRollerHandler.NewWithRoomChecker(services.DiceRoller, services.Room)
	} else {
		diceHandler = diceRollerHandler.New(services.DiceRoller)
	}

	handler := &Handler{
		domainHandlers: &domainHandlers{
			characterHandler:  characterHandler.New(services.Character),
			diceRollerHandler: diceHandler,
			roomHandler:       roomHandler.NewWithHub(services.Room, roomHub),
			userHandler:       userHandler.New(services.User),
		},
		portraitFileServer: dependencies.PortraitFileServer,
	}
	return handler
}

func (h *Handler) RoomHub() *ws.RoomHub {
	return h.domainHandlers.roomHandler.Hub()
}

func (h *Handler) CloseDeletedRoomSockets(result roomModel.CleanupRoomsResult, reason string) {
	h.domainHandlers.roomHandler.CloseDeletedRooms(result.DeletedRoomIDs, reason)
}
