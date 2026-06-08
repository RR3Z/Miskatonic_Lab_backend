package tests

import (
	"context"
	"sync"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
)

type testEvent struct {
	name string
}

type otherTestEvent struct {
	name string
}

func (e testEvent) EventName() string {
	return e.name
}

func (e otherTestEvent) EventName() string {
	return e.name
}

type eventSubscriptionCase struct {
	name             string
	subscribedEvent  events.Event
	publishedEvent   events.Event
	expectedRecorded events.Event
}

func eventSubscriptionCases() []eventSubscriptionCase {
	return []eventSubscriptionCase{
		{
			name:             "subscribe value publish pointer",
			subscribedEvent:  testEvent{},
			publishedEvent:   &testEvent{name: "test.event"},
			expectedRecorded: &testEvent{name: "test.event"},
		},
		{
			name:             "subscribe pointer publish value",
			subscribedEvent:  &testEvent{},
			publishedEvent:   testEvent{name: "test.event"},
			expectedRecorded: testEvent{name: "test.event"},
		},
	}
}

type recordingEventHandler struct {
	mu     sync.Mutex
	calls  int
	events []events.Event
	name   string
	order  *[]string
}

func (h *recordingEventHandler) Handle(_ context.Context, event events.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, event)
	h.calls++
	if h.order != nil {
		*h.order = append(*h.order, h.name)
	}
}

func (h *recordingEventHandler) Calls() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.calls
}

func (h *recordingEventHandler) Events() []events.Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	return append([]events.Event(nil), h.events...)
}
