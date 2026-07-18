package character

import (
	"context"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type ICharacter interface {
	GetAllCharacters(ctx context.Context, userID string) ([]characterDTO.CharacterSummaryModel, error)
	GetCharacter(ctx context.Context, input characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error)
	CreateCharacter(ctx context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error)
	UpdateCharacter(ctx context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error)
	PatchCharacter(ctx context.Context, input characterDTO.PatchCharacterInput) (characterDTO.CharacterShortModel, error)
	ReplacePortrait(ctx context.Context, input characterDTO.ReplacePortraitInput) (characterDTO.CharacterShortModel, error)
	DeleteCharacter(ctx context.Context, input characterDTO.DeleteCharacterInput) error

	GetHealth(ctx context.Context, input healthDTO.GetHealthInput) (db.HealthState, error)
	UpsertHealth(ctx context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error)
	DeleteHealth(ctx context.Context, input healthDTO.DeleteHealthInput) error

	GetSanity(ctx context.Context, input sanityDTO.GetSanityInput) (db.SanityState, error)
	UpsertSanity(ctx context.Context, input sanityDTO.UpsertSanityInput) (db.SanityState, error)
	DeleteSanity(ctx context.Context, input sanityDTO.DeleteSanityInput) error

	GetMagic(ctx context.Context, input magicDTO.GetMagicInput) (db.MagicState, error)
	UpsertMagic(ctx context.Context, input magicDTO.UpsertMagicInput) (db.MagicState, error)
	DeleteMagic(ctx context.Context, input magicDTO.DeleteMagicInput) error

	GetLuck(ctx context.Context, input luckDTO.GetLuckInput) (db.LuckState, error)
	UpsertLuck(ctx context.Context, input luckDTO.UpsertLuckInput) (db.LuckState, error)
	DeleteLuck(ctx context.Context, input luckDTO.DeleteLuckInput) error

	GetFinances(ctx context.Context, input financesDTO.GetFinancesInput) (db.Finance, error)
	UpsertFinances(ctx context.Context, input financesDTO.UpsertFinancesInput) (db.Finance, error)
	DeleteFinances(ctx context.Context, input financesDTO.DeleteFinancesInput) error

	GetInventoryItems(ctx context.Context, input inventoryDTO.GetInventoryItemsInput) ([]db.CharacterInventoryItem, error)
	GetInventoryItem(ctx context.Context, input inventoryDTO.GetInventoryItemInput) (db.CharacterInventoryItem, error)
	CreateInventoryItem(ctx context.Context, input inventoryDTO.CreateInventoryItemInput) (db.CharacterInventoryItem, error)
	UpdateInventoryItem(ctx context.Context, input inventoryDTO.UpdateInventoryItemInput) (db.CharacterInventoryItem, error)
	DeleteInventoryItem(ctx context.Context, input inventoryDTO.DeleteInventoryItemInput) error

	GetBackstory(ctx context.Context, input backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error)
	UpsertBackstory(ctx context.Context, input backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error)
	DeleteBackstory(ctx context.Context, input backstoriesDTO.DeleteBackstoryInput) error
	GetBackstoryItems(ctx context.Context, input backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error)
	GetBackstoryItem(ctx context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error)
	CreateBackstoryItem(ctx context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error)
	UpdateBackstoryItem(ctx context.Context, input backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error)
	DeleteBackstoryItem(ctx context.Context, input backstoriesDTO.DeleteBackstoryItemInput) error

	GetSkills(ctx context.Context, input skillsDTO.GetSkillsInput) ([]skillsDTO.SkillModel, error)
	GetSkill(ctx context.Context, input skillsDTO.GetSkillInput) (skillsDTO.SkillModel, error)
	CreateSkill(ctx context.Context, input skillsDTO.CreateSkillInput) (skillsDTO.SkillModel, error)
	UpdateSkill(ctx context.Context, input skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error)
	DeleteSkill(ctx context.Context, input skillsDTO.DeleteSkillInput) error

	GetDerivedStats(ctx context.Context, input derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error)

	GetCharacteristics(ctx context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error)
	UpsertCharacteristics(ctx context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error)
	DeleteCharacteristics(ctx context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error

	GetNotes(ctx context.Context, input notesDTO.GetNotesInput) ([]db.Note, error)
	GetNote(ctx context.Context, input notesDTO.GetNoteInput) (db.Note, error)
	CreateNote(ctx context.Context, input notesDTO.CreateNoteInput) (db.Note, error)
	UpdateNote(ctx context.Context, input notesDTO.UpdateNoteInput) (db.Note, error)
	DeleteNote(ctx context.Context, input notesDTO.DeleteNoteInput) error
}
