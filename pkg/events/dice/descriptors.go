package dice

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"

const domain = "dice"

var descriptors = []events.EventDescriptor{
	descriptor(DiceRollMakeSucceeded{}, "dice_roll", "make", events.OutcomeSucceeded),
	descriptor(DiceRollMakeFailed{}, "dice_roll", "make", events.OutcomeFailed),
	descriptor(DiceRollsListSucceeded{}, "dice_rolls", "list", events.OutcomeSucceeded),
	descriptor(DiceRollsListFailed{}, "dice_rolls", "list", events.OutcomeFailed),
}

var roomPublishingDescriptors = []events.EventDescriptor{
	descriptor(DiceRollMakeSucceeded{}, "dice_roll", "make", events.OutcomeSucceeded),
}

func Descriptors() []events.EventDescriptor {
	return append([]events.EventDescriptor(nil), descriptors...)
}

func RoomPublishingEvents() []events.Event {
	return events.EventPrototypes(roomPublishingDescriptors)
}

func descriptor(event events.Event, resource string, action string, outcome events.EventOutcome) events.EventDescriptor {
	return events.EventDescriptor{
		Event:    event,
		Domain:   domain,
		Resource: resource,
		Action:   action,
		Outcome:  outcome,
	}
}
