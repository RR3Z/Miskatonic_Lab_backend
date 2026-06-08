package publishers

import (
	"context"
	"reflect"

	e "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
)

// Implements EventSubscriber and EventPublisher (EventSubscriberPublisher)
type SyncPublisher struct {
	allHandlers []e.EventHandler
	handlers    map[reflect.Type][]e.EventHandler
}

func NewSyncPublisher() *SyncPublisher {
	return &SyncPublisher{
		allHandlers: []e.EventHandler{},
		handlers:    map[reflect.Type][]e.EventHandler{},
	}
}

func (p *SyncPublisher) Subscribe(handler e.EventHandler) {
	p.allHandlers = append(p.allHandlers, handler)
}

func (p *SyncPublisher) SubscribeEvent(event e.Event, handler e.EventHandler) {
	t := e.EventTypeOf(event)
	p.handlers[t] = append(p.handlers[t], handler)
}

func (p *SyncPublisher) Publish(ctx context.Context, event e.Event) {
	for _, handler := range p.allHandlers {
		handler.Handle(ctx, event)
	}

	for _, handler := range p.handlers[e.EventTypeOf(event)] {
		handler.Handle(ctx, event)
	}
}
