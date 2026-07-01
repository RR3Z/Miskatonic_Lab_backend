package character

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	characterModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type ICharacter interface {
	GetAllCharacters(ctx context.Context, userID string) ([]characterModel.CharacterShortModel, error)
	GetCharacter(ctx context.Context, input characterModel.GetCharacterInput) (characterModel.CharacterModel, error)
	CreateCharacter(ctx context.Context, input characterModel.CreateCharacterInput) (characterModel.CharacterShortModel, error)
	UpdateCharacter(ctx context.Context, input characterModel.UpdateCharacterInput) (characterModel.CharacterShortModel, error)
	DeleteCharacter(ctx context.Context, input characterModel.DeleteCharacterInput) error

	GetHealth(ctx context.Context, input characterModel.GetHealthInput) (db.HealthState, error)
	UpsertHealth(ctx context.Context, input characterModel.UpsertHealthInput) (db.HealthState, error)
	DeleteHealth(ctx context.Context, input characterModel.DeleteHealthInput) error

	GetSanity(ctx context.Context, input characterModel.GetSanityInput) (db.SanityState, error)
	UpsertSanity(ctx context.Context, input characterModel.UpsertSanityInput) (db.SanityState, error)
	DeleteSanity(ctx context.Context, input characterModel.DeleteSanityInput) error

	GetMagic(ctx context.Context, input characterModel.GetMagicInput) (db.MagicState, error)
	UpsertMagic(ctx context.Context, input characterModel.UpsertMagicInput) (db.MagicState, error)
	DeleteMagic(ctx context.Context, input characterModel.DeleteMagicInput) error

	GetLuck(ctx context.Context, input characterModel.GetLuckInput) (db.LuckState, error)
	UpsertLuck(ctx context.Context, input characterModel.UpsertLuckInput) (db.LuckState, error)
	DeleteLuck(ctx context.Context, input characterModel.DeleteLuckInput) error

	GetFinances(ctx context.Context, input characterModel.GetFinancesInput) (db.Finance, error)
	UpsertFinances(ctx context.Context, input characterModel.UpsertFinancesInput) (db.Finance, error)
	DeleteFinances(ctx context.Context, input characterModel.DeleteFinancesInput) error

	GetBackstory(ctx context.Context, input characterModel.GetBackstoryInput) (model.BackstoryModel, error)
	UpsertBackstory(ctx context.Context, input characterModel.UpsertBackstoryInput) (model.BackstoryModel, error)
	DeleteBackstory(ctx context.Context, input characterModel.DeleteBackstoryInput) error
	GetBackstoryItems(ctx context.Context, input characterModel.GetBackstoryItemsInput) ([]model.BackstoryItemModel, error)
	GetBackstoryItem(ctx context.Context, input characterModel.GetBackstoryItemInput) (model.BackstoryItemModel, error)
	CreateBackstoryItem(ctx context.Context, input characterModel.CreateBackstoryItemInput) (model.BackstoryItemModel, error)
	UpdateBackstoryItem(ctx context.Context, input characterModel.UpdateBackstoryItemInput) (model.BackstoryItemModel, error)
	DeleteBackstoryItem(ctx context.Context, input characterModel.DeleteBackstoryItemInput) error

	GetSkills(ctx context.Context, input characterModel.GetSkillsInput) ([]model.SkillModel, error)
	GetSkill(ctx context.Context, input characterModel.GetSkillInput) (model.SkillModel, error)
	CreateSkill(ctx context.Context, input characterModel.CreateSkillInput) (model.SkillModel, error)
	UpdateSkill(ctx context.Context, input characterModel.UpdateSkillInput) (model.SkillModel, error)
	DeleteSkill(ctx context.Context, input characterModel.DeleteSkillInput) error

	GetDerivedStats(ctx context.Context, input characterModel.GetDerivedStatsInput) (db.DerivedStat, error)
	UpsertDerivedStats(ctx context.Context, input characterModel.UpsertDerivedStatsInput) (db.DerivedStat, error)
	DeleteDerivedStats(ctx context.Context, input characterModel.DeleteDerivedStatsInput) error

	GetCharacteristics(ctx context.Context, input characterModel.GetCharacteristicsInput) (db.Characteristic, error)
	UpsertCharacteristics(ctx context.Context, input characterModel.UpsertCharacteristicsInput) (db.Characteristic, error)
	DeleteCharacteristics(ctx context.Context, input characterModel.DeleteCharacteristicsInput) error

	GetNotes(ctx context.Context, input characterModel.GetNotesInput) ([]db.Note, error)
	GetNote(ctx context.Context, input characterModel.GetNoteInput) (db.Note, error)
	CreateNote(ctx context.Context, input characterModel.CreateNoteInput) (db.Note, error)
	UpdateNote(ctx context.Context, input characterModel.UpdateNoteInput) (db.Note, error)
	DeleteNote(ctx context.Context, input characterModel.DeleteNoteInput) error
}
