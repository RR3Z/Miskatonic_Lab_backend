package roomDTO

import (
	"encoding/json"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomEventModel struct {
	ID        pgtype.UUID        `json:"id"`
	RoomID    pgtype.UUID        `json:"room_id"`
	ActorID   string             `json:"actor_id"`
	Type      string             `json:"type"`
	Payload   json.RawMessage    `json:"payload"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func ToRoomEventModel(event db.RoomEvent) RoomEventModel {
	payload := json.RawMessage(event.Payload)
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}

	return RoomEventModel{
		ID:        event.ID,
		RoomID:    event.RoomID,
		ActorID:   event.ActorID,
		Type:      event.EventType,
		Payload:   payload,
		CreatedAt: event.CreatedAt,
	}
}

func ToRoomEventModels(events []db.RoomEvent) []RoomEventModel {
	models := make([]RoomEventModel, len(events))
	for i, event := range events {
		models[i] = ToRoomEventModel(event)
	}
	return models
}
