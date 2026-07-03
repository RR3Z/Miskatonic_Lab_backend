package events

import "reflect"

type EventOutcome string

const (
	OutcomeSucceeded EventOutcome = "succeeded"
	OutcomeFailed    EventOutcome = "failed"
	OutcomeSkipped   EventOutcome = "skipped"
)

type EventDescriptor struct {
	Event    Event
	Domain   string
	Resource string
	Action   string
	Outcome  EventOutcome
}

type DescriptorRegistry struct {
	descriptors map[reflect.Type]EventDescriptor
}

func NewDescriptorRegistry(groups ...[]EventDescriptor) DescriptorRegistry {
	descriptors := make(map[reflect.Type]EventDescriptor)
	for _, group := range groups {
		for _, descriptor := range group {
			descriptors[EventTypeOf(descriptor.Event)] = descriptor
		}
	}

	return DescriptorRegistry{descriptors: descriptors}
}

func (r DescriptorRegistry) Describe(event Event) (EventDescriptor, bool) {
	if event == nil {
		return EventDescriptor{}, false
	}

	descriptor, ok := r.descriptors[EventTypeOf(event)]
	return descriptor, ok
}

func EventPrototypes(descriptors []EventDescriptor) []Event {
	events := make([]Event, len(descriptors))
	for i, descriptor := range descriptors {
		events[i] = descriptor.Event
	}
	return events
}
