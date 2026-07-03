package main

import (
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	roomListeners "github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners/room"
	EventsLogging "github.com/RR3Z/Miskatonic_Lab_backend/pkg/observability/logging"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
)

func registerEventListeners(eventBus *events.EventBus, services *appService.Service, appHandlers *handler.Handler) {
	eventBus.SubscribeAllSync(EventsLogging.NewDefaultEventLogger(slog.Default()))

	characterRoomListener := roomListeners.NewCharacterRoomListener(services.Room, appHandlers.RoomHub())
	for _, event := range roomListeners.MutationCharacterEvents() {
		eventBus.SubscribeAsync(event, characterRoomListener)
	}

	diceRoomListener := roomListeners.NewDiceRollerRoomListener(services.Room, appHandlers.RoomHub())
	for _, event := range diceEvents.RoomPublishingEvents() {
		eventBus.SubscribeAsync(event, diceRoomListener)
	}
}
