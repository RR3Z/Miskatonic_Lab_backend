package main

import (
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	roomListeners "github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners/room"
	EventsLogging "github.com/RR3Z/Miskatonic_Lab_backend/pkg/observability/logging"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
)

func registerEventListeners(eventBus *events.EventBus, services *appService.Service, appHandlers *handler.Handler) {
	registerEventLogging(eventBus)
	registerRoomEventListeners(eventBus, services, appHandlers)
}

func registerEventLogging(eventBus *events.EventBus) {
	eventBus.SubscribeAllSync(EventsLogging.NewDefaultEventLogger(slog.Default()))
}

func registerRoomEventListeners(eventBus *events.EventBus, services *appService.Service, appHandlers *handler.Handler) {
	characterRoomListener := roomListeners.NewCharacterRoomListener(services.Room, appHandlers.RoomHub())
	subscribeAsyncEvents(eventBus, characterEvents.RoomMutationEvents(), characterRoomListener)

	diceRoomListener := roomListeners.NewDiceRollerRoomListener(services.Room, appHandlers.RoomHub())
	subscribeAsyncEvents(eventBus, diceEvents.RoomPublishingEvents(), diceRoomListener)
}

func subscribeAsyncEvents(eventBus *events.EventBus, eventPrototypes []events.Event, handler events.EventHandler) {
	for _, event := range eventPrototypes {
		eventBus.SubscribeAsync(event, handler)
	}
}
