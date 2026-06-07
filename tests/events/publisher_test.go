package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/stretchr/testify/require"
)

type testEvent struct {
	name string
}

func (e testEvent) EventName() string {
	return e.name
}

type recordingEventHandler struct {
	calls  int
	events []events.Event
	name   string
	order  *[]string
}

func (h *recordingEventHandler) Handle(_ context.Context, event events.Event) {
	h.calls++
	h.events = append(h.events, event)
	if h.order != nil {
		*h.order = append(*h.order, h.name)
	}
}

func TestSyncPublisherPublishWithNoHandlersDoesNotPanic(t *testing.T) {
	publisher := events.NewSyncPublisher()

	require.NotPanics(t, func() {
		publisher.Publish(context.Background(), testEvent{name: "test.event"})
	})
}

func TestSyncPublisherPublishesEventToSubscribedHandlersInOrder(t *testing.T) {
	publisher := events.NewSyncPublisher()
	order := []string{}
	first := &recordingEventHandler{name: "first", order: &order}
	second := &recordingEventHandler{name: "second", order: &order}
	event := testEvent{name: "test.event"}

	publisher.Subscribe(first)
	publisher.Subscribe(second)
	publisher.Publish(context.Background(), event)

	require.Equal(t, 1, first.calls)
	require.Equal(t, 1, second.calls)
	require.Equal(t, []events.Event{event}, first.events)
	require.Equal(t, []events.Event{event}, second.events)
	require.Equal(t, []string{"first", "second"}, order)
}

func TestNoopPublisherDoesNotPanic(t *testing.T) {
	publisher := &events.NoopPublisher{}

	require.NotPanics(t, func() {
		publisher.Publish(context.Background(), testEvent{name: "test.event"})
	})
}
