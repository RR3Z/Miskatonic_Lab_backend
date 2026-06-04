package events

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event Event)
}
