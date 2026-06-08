package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/publishers"
	"github.com/stretchr/testify/require"
)

func TestSyncPublisherPublishWithNoHandlersDoesNotPanic(t *testing.T) {
	publisher := publishers.NewSyncPublisher()

	require.NotPanics(t, func() {
		publisher.Publish(context.Background(), testEvent{name: "test.event"})
	})
}

func TestSyncPublisherPublishesEventToSubscribedHandlersInOrder(t *testing.T) {
	publisher := publishers.NewSyncPublisher()
	order := []string{}
	first := &recordingEventHandler{name: "first", order: &order}
	second := &recordingEventHandler{name: "second", order: &order}
	event := testEvent{name: "test.event"}

	publisher.Subscribe(first)
	publisher.Subscribe(second)
	publisher.Publish(context.Background(), event)

	require.Equal(t, 1, first.Calls())
	require.Equal(t, 1, second.Calls())
	require.Equal(t, []events.Event{event}, first.Events())
	require.Equal(t, []events.Event{event}, second.Events())
	require.Equal(t, []string{"first", "second"}, order)
}

func TestSyncPublisherPublishesEventOnlyToMatchingEventHandlers(t *testing.T) {
	publisher := publishers.NewSyncPublisher()
	matching := &recordingEventHandler{name: "matching"}
	other := &recordingEventHandler{name: "other"}
	event := testEvent{name: "test.event"}

	publisher.SubscribeEvent(testEvent{}, matching)
	publisher.SubscribeEvent(otherTestEvent{}, other)
	publisher.Publish(context.Background(), event)

	require.Equal(t, 1, matching.Calls())
	require.Equal(t, []events.Event{event}, matching.Events())
	require.Equal(t, 0, other.Calls())
}

func TestSyncPublisherPublishesToAllHandlersBeforeEventHandlers(t *testing.T) {
	publisher := publishers.NewSyncPublisher()
	order := []string{}
	all := &recordingEventHandler{name: "all", order: &order}
	typed := &recordingEventHandler{name: "typed", order: &order}
	event := testEvent{name: "test.event"}

	publisher.Subscribe(all)
	publisher.SubscribeEvent(testEvent{}, typed)
	publisher.Publish(context.Background(), event)

	require.Equal(t, 1, all.Calls())
	require.Equal(t, 1, typed.Calls())
	require.Equal(t, []string{"all", "typed"}, order)
}

func TestSyncPublisherMatchesEventSubscriptionByValueAndPointerType(t *testing.T) {
	for _, tc := range eventSubscriptionCases() {
		t.Run(tc.name, func(t *testing.T) {
			publisher := publishers.NewSyncPublisher()
			handler := &recordingEventHandler{name: "handler"}

			publisher.SubscribeEvent(tc.subscribedEvent, handler)
			publisher.Publish(context.Background(), tc.publishedEvent)

			require.Equal(t, 1, handler.Calls())
			require.Equal(t, []events.Event{tc.expectedRecorded}, handler.Events())
		})
	}
}
