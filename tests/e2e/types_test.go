package tests

import (
	"encoding/json"
	"net/http"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type e2eSubject struct {
	baseURL string
	token   string
	userID  string
	client  *http.Client
	pool    *pgxpool.Pool
	queries *db.Queries
}

type e2eUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type e2eIDResponse struct {
	ID string `json:"id"`
}

type e2eCharacterResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	HP     struct {
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
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
