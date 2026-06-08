package events

import "context"

type EventBus struct {
	syncPublisher  EventSubscriberPublisher
	asyncPublisher EventSubscriberPublisher
}

func NewEventBus(syncPublisher EventSubscriberPublisher, asyncPublisher EventSubscriberPublisher) *EventBus {
	return &EventBus{
		syncPublisher:  syncPublisher,
		asyncPublisher: asyncPublisher,
	}
}

func (b *EventBus) Publish(ctx context.Context, event Event) {
	b.syncPublisher.Publish(ctx, event)
	b.asyncPublisher.Publish(ctx, event)
}

func (b *EventBus) SubscribeSync(event Event, handler EventHandler) {
	b.syncPublisher.SubscribeEvent(event, handler)
}

func (b *EventBus) SubscribeAsync(event Event, handler EventHandler) {
	b.asyncPublisher.SubscribeEvent(event, handler)
}

func (b *EventBus) SubscribeAllSync(handler EventHandler) {
	b.syncPublisher.Subscribe(handler)
}

func (b *EventBus) SubscribeAllAsync(handler EventHandler) {
	b.asyncPublisher.Subscribe(handler)
}
