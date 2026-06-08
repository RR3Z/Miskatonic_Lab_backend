package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/publishers"
	"github.com/stretchr/testify/require"
)

func TestAsyncPublisherPublishWithNoHandlersDoesNotPanic(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(1, nil)

	require.NotPanics(t, func() {
		publisher.Publish(context.Background(), testEvent{name: "test.event"})
	})
}

func TestAsyncPublisherDoesNotHandleQueuedEventsBeforeStart(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(1, nil)
	handler := &recordingEventHandler{name: "handler"}

	publisher.SubscribeEvent(testEvent{}, handler)
	publisher.Publish(context.Background(), testEvent{name: "test.event"})

	time.Sleep(20 * time.Millisecond)
	require.Equal(t, 0, handler.Calls())
}

func TestAsyncPublisherPublishesEventOnlyToMatchingEventHandlers(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(10, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher.Start(ctx, 1)

	matching := &recordingEventHandler{name: "matching"}
	other := &recordingEventHandler{name: "other"}
	event := testEvent{name: "test.event"}

	publisher.SubscribeEvent(testEvent{}, matching)
	publisher.SubscribeEvent(otherTestEvent{}, other)
	publisher.Publish(context.Background(), event)

	require.Eventually(t, func() bool {
		return matching.Calls() == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, []events.Event{event}, matching.Events())
	require.Equal(t, 0, other.Calls())
}

func TestAsyncPublisherPublishesToAllAndEventHandlers(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(10, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher.Start(ctx, 1)

	all := &recordingEventHandler{name: "all"}
	typed := &recordingEventHandler{name: "typed"}
	event := testEvent{name: "test.event"}

	publisher.Subscribe(all)
	publisher.SubscribeEvent(testEvent{}, typed)
	publisher.Publish(context.Background(), event)

	require.Eventually(t, func() bool {
		return all.Calls() == 1 && typed.Calls() == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, []events.Event{event}, all.Events())
	require.Equal(t, []events.Event{event}, typed.Events())
}

func TestAsyncPublisherMatchesEventSubscriptionByValueAndPointerType(t *testing.T) {
	for _, tc := range eventSubscriptionCases() {
		t.Run(tc.name, func(t *testing.T) {
			publisher := publishers.NewAsyncPublisher(10, nil)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			publisher.Start(ctx, 1)

			handler := &recordingEventHandler{name: "handler"}

			publisher.SubscribeEvent(tc.subscribedEvent, handler)
			publisher.Publish(context.Background(), tc.publishedEvent)

			require.Eventually(t, func() bool {
				return handler.Calls() == 1
			}, time.Second, 10*time.Millisecond)
			require.Equal(t, []events.Event{tc.expectedRecorded}, handler.Events())
		})
	}
}

func TestAsyncPublisherDropsEventWhenQueueIsFull(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(0, nil)
	handler := &recordingEventHandler{name: "handler"}

	publisher.SubscribeEvent(testEvent{}, handler)
	require.NotPanics(t, func() {
		publisher.Publish(context.Background(), testEvent{name: "test.event"})
	})

	time.Sleep(20 * time.Millisecond)
	require.Equal(t, 0, handler.Calls())
}

func TestAsyncPublisherStartWithZeroWorkersDoesNotHandleEvents(t *testing.T) {
	publisher := publishers.NewAsyncPublisher(1, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher.Start(ctx, 0)
	handler := &recordingEventHandler{name: "handler"}

	publisher.SubscribeEvent(testEvent{}, handler)
	publisher.Publish(context.Background(), testEvent{name: "test.event"})

	time.Sleep(20 * time.Millisecond)
	require.Equal(t, 0, handler.Calls())
}
