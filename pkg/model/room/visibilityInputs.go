package roomDTO

import "github.com/jackc/pgx/v5/pgtype"

type ListSelectedCharactersInput struct {
	RoomID pgtype.UUID
	UserID string
}
