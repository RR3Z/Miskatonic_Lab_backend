package roomEvents

type EventType string

const (
	EventChatMessage      EventType = "chat.message"
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
	CharacterID string `json:"character_id"`
	Expression  string `json:"expression"`
	Result      int32  `json:"result"`
}

type OwnerTransferredPayload struct {
	PreviousOwnerID string `json:"previous_owner_id"`
	NewOwnerID      string `json:"new_owner_id"`
}
