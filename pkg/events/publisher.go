package events

import (
	"context"
	"reflect"
)

type EventPublisher interface {
	Publish(ctx context.Context, event Event)
}

type EventSubscriber interface {
	Subscribe(handler EventHandler)
	SubscribeEvent(event Event, handler EventHandler)
}

type EventSubscriberPublisher interface {
	EventPublisher
	EventSubscriber
}

func EventTypeOf(event Event) reflect.Type {
	t := reflect.TypeOf(event)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}
