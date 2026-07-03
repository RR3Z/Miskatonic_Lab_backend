package main

import (
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners"
	listenerHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners/helpers"
	EventsLogging "github.com/RR3Z/Miskatonic_Lab_backend/pkg/observability/logging"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
)

func registerEventListeners(eventBus *events.EventBus, services *appService.Service, appHandlers *handler.Handler) {
	eventBus.SubscribeAllSync(EventsLogging.NewCharacterEventLogger(slog.Default()))

	characterRoomListener := listeners.NewCharacterRoomListener(services.Room, appHandlers.RoomHub())
	for _, event := range listenerHelpers.MutationCharacterEvents() {
		eventBus.SubscribeAsync(event, characterRoomListener)
	}

	eventBus.SubscribeAsync(
		diceEvents.DiceRollMakeSucceeded{},
		listeners.NewDiceRollerRoomListener(services.Room, appHandlers.RoomHub()),
	)
}
