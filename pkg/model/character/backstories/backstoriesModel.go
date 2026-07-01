package backstoriesDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type BackstoryItemModel struct {
	ID pgtype.UUID `json:"id"`

	Section string `json:"section"`
	Title   string `json:"title"`
	Text    string `json:"text"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func ToBackstoryItemModel(item db.BackstoryItem) BackstoryItemModel {
	return BackstoryItemModel{
		ID:        item.ID,
		Section:   item.Section,
		Title:     item.Title,
		Text:      item.Text,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func ToBackstoryItemModels(items []db.BackstoryItem) []BackstoryItemModel {
	models := make([]BackstoryItemModel, len(items))
	for i, item := range items {
		models[i] = ToBackstoryItemModel(item)
	}
	return models
}

type BackstoryModel struct {
	ID          pgtype.UUID `json:"id"`
	CharacterID pgtype.UUID `json:"character_id"`

	PersonalDescription *string              `json:"personal_description"`
	Items               []BackstoryItemModel `json:"items"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func ToBackstoryModel(b db.Backstory, items []db.BackstoryItem) BackstoryModel {
	return BackstoryModel{
		ID:                  b.ID,
		CharacterID:         b.CharacterID,
		PersonalDescription: b.PersonalDescription,
		Items:               ToBackstoryItemModels(items),
		CreatedAt:           b.CreatedAt,
		UpdatedAt:           b.UpdatedAt,
	}
}
