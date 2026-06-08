package publishers

import (
	"context"
	"log/slog"
	"reflect"

	e "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
)

type asyncEventJob struct {
	event   e.Event
	handler e.EventHandler
}

// Implements EventSubscriber and EventPublisher (EventSubscriberPublisher)
type AsyncPublisher struct {
	allHandlers []e.EventHandler
	handlers    map[reflect.Type][]e.EventHandler
	jobs        chan asyncEventJob
	logger      *slog.Logger
}

func NewAsyncPublisher(bufferSize int, logger *slog.Logger) *AsyncPublisher {
	return &AsyncPublisher{
		allHandlers: []e.EventHandler{},
		handlers:    map[reflect.Type][]e.EventHandler{},
		jobs:        make(chan asyncEventJob, bufferSize),
		logger:      logger,
	}
}

func (p *AsyncPublisher) Subscribe(handler e.EventHandler) {
	p.allHandlers = append(p.allHandlers, handler)
}

func (p *AsyncPublisher) SubscribeEvent(event e.Event, handler e.EventHandler) {
	t := e.EventTypeOf(event)
	p.handlers[t] = append(p.handlers[t], handler)
}

func (p *AsyncPublisher) Publish(ctx context.Context, event e.Event) {
	p.publishToHandlers(ctx, event, p.allHandlers)
	p.publishToHandlers(ctx, event, p.handlers[e.EventTypeOf(event)])
}

func (p *AsyncPublisher) Start(ctx context.Context, workersCount int) {
	for range workersCount {
		go p.work(ctx)
	}
}

func (p *AsyncPublisher) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobs:
			job.handler.Handle(context.Background(), job.event)
		}
	}
}

func (p *AsyncPublisher) publishToHandlers(_ context.Context, event e.Event, handlers []e.EventHandler) {
	for _, handler := range handlers {
		select {
		case p.jobs <- asyncEventJob{event: event, handler: handler}:
		default:
			if p.logger != nil {
				p.logger.Warn(
					"async event queue is full",
					"event", event.EventName(),
				)
			}
		}
	}
}
