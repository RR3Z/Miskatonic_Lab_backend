package roomEvents

import "encoding/json"

type EventType string

const (
	EventChatMessage      EventType = "chat.message"
	EventCharacterChanged EventType = "character.changed"
	EventCommandError     EventType = "command.error"
	EventDiceRoll         EventType = "dice.roll"
	EventMemberJoined     EventType = "member.joined"
	EventMemberLeft       EventType = "member.left"
	EventOwnerTransferred EventType = "owner.transferred"
)

type Event struct {
	Type    string `json:"type"`
	RoomID  string `json:"room_id"`
	ActorID string `json:"actor_id"`
	Payload any    `json:"payload,omitempty"`
}

type ChatMessagePayload struct {
	Text string `json:"text"`
}

type DiceRollPayload struct {
	RollID      string          `json:"roll_id"`
	CharacterID string          `json:"character_id"`
	Expression  string          `json:"expression"`
	Result      int32           `json:"result"`
	Details     json.RawMessage `json:"details"`
}

type OwnerTransferredPayload struct {
	PreviousOwnerID string `json:"previous_owner_id"`
	NewOwnerID      string `json:"new_owner_id"`
}

type CharacterChangedPayload struct {
	CharacterID string  `json:"character_id"`
	Resource    string  `json:"resource"`
	Action      string  `json:"action"`
	ResourceID  *string `json:"resource_id,omitempty"`
	SourceEvent *string `json:"source_event,omitempty"`
}
