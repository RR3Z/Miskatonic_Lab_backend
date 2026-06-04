package events

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, event Event)
}

type NoopPublisher struct{}

func (p *NoopPublisher) Publish(ctx context.Context, event Event) {}

type SyncPublisher struct {
	handlers []EventHandler
}

func NewSyncPublisher() *SyncPublisher {
	return &SyncPublisher{}
}

func (p *SyncPublisher) Subscribe(handler EventHandler) {
	p.handlers = append(p.handlers, handler)
}

func (p *SyncPublisher) Publish(ctx context.Context, event Event) {
	for _, handler := range p.handlers {
		handler.Handle(ctx, event)
	}
}
