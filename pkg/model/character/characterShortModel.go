package character

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterShortModel struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`

	Name       string  `json:"name"`
	PlayerName *string `json:"player_name"`
	Occupation *string `json:"occupation"`
	Age        *int16  `json:"age"`
	Sex        *string `json:"sex"`
	Residence  *string `json:"residence"`
	Birthplace *string `json:"birthplace"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func ToCharacterShortModel(c db.Character) CharacterShortModel {
	return CharacterShortModel{
		ID:         c.ID,
		UserID:     c.UserID,
		Name:       c.Name,
		PlayerName: c.PlayerName,
		Occupation: c.Occupation,
		Age:        c.Age,
		Sex:        c.Sex,
		Residence:  c.Residence,
		Birthplace: c.Birthplace,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}
