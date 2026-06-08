package publishers

import (
	"context"

	e "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
)

// Implements EventPublisher.
type NoopPublisher struct{}

func (p *NoopPublisher) Publish(ctx context.Context, event e.Event) {}

func (p *NoopPublisher) Subscribe(handler e.EventHandler) {}

func (p *NoopPublisher) SubscribeEvent(event e.Event, handler e.EventHandler) {}
