package character

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"

// Implements ICharacter
type EventPublishingCharacterService struct {
	next      ICharacter
	publisher events.EventPublisher
}

func NewEventPublishingCharacterService(next ICharacter, publisher events.EventPublisher) *EventPublishingCharacterService {
	return &EventPublishingCharacterService{
		next:      next,
		publisher: publisher,
	}
}
