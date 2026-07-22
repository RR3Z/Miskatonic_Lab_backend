package tests

import (
	"encoding/json"
	"net/http"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type e2eSubject struct {
	baseURL  string
	identity e2eClerkIdentity
	userID   string
	client   *http.Client
	pool     *pgxpool.Pool
	queries  *db.Queries
}

type e2eUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type e2eIDResponse struct {
	ID string `json:"id"`
}

type e2eCharacterShortResponse struct {
	ID          string  `json:"id"`
	PortraitURL *string `json:"portrait_url"`
}

type e2eCharacterSummaryResponse struct {
	ID          string  `json:"id"`
	PortraitURL *string `json:"portrait_url"`
}

type e2eErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type e2eCharacterResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	PortraitURL *string `json:"portrait_url"`
	HP          struct {
		MaxHp     int16 `json:"max_hp"`
		CurrentHp int16 `json:"current_hp"`
	} `json:"hp"`
}

type e2eHealthResponse struct {
	MaxHp     int16 `json:"max_hp"`
	CurrentHp int16 `json:"current_hp"`
}

type e2eRoomResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
}

type e2eRoomEventResponse struct {
	Sequence int64           `json:"sequence"`
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

type e2eSelectedCharacterResponse struct {
	UserID    string               `json:"user_id"`
	Role      string               `json:"role"`
	Character e2eCharacterResponse `json:"character"`
}

type e2eRoomCommand struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type e2eRoomSocketEvent struct {
	Type     string          `json:"type"`
	RoomID   string          `json:"room_id"`
	Sequence int64           `json:"sequence"`
	ActorID  string          `json:"actor_id"`
	Payload  json.RawMessage `json:"payload"`
}

type e2eChatPayload struct {
	Text string `json:"text"`
}
