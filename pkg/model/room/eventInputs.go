package roomDTO

import "github.com/jackc/pgx/v5/pgtype"

type CreateChatMessageInput struct {
	RoomID  pgtype.UUID
	ActorID string
	Text    string
}

type CreateDiceRollRoomEventInput struct {
	RoomID      pgtype.UUID
	ActorID     string
	RollID      string
	CharacterID string
	Expression  string
	Result      int32
	Details     []byte
}

type CharacterChangedRoomEventChange struct {
	Resource    string
	Action      string
	ResourceID  *string
	SourceEvent *string
}

type CreateCharacterChangedRoomEventsInput struct {
	CharacterID pgtype.UUID
	ActorID     string
	Change      CharacterChangedRoomEventChange
}
