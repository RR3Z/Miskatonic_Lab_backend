package character

import "github.com/jackc/pgx/v5/pgtype"

type GetCharacterInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type CreateCharacterInput struct {
	UserID     string
	Name       string
	PlayerName *string
	Occupation *string
	Age        *int16
	Sex        *string
	Residence  *string
	Birthplace *string
}

type UpdateCharacterInput struct {
	UserID     string
	ID         pgtype.UUID
	Name       string
	PlayerName *string
	Occupation *string
	Age        *int16
	Sex        *string
	Residence  *string
	Birthplace *string
}

type DeleteCharacterInput struct {
	UserID string
	ID     pgtype.UUID
}

type GetHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
	MaxHp       *int16
	CurrentHp   *int16
	MajorWound  *bool
	Unconscious *bool
	Dying       *bool
	Dead        *bool
}

type DeleteHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetSanityInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertSanityInput struct {
	UserID        string
	CharacterID   pgtype.UUID
	MaxSanity     *int16
	CurrentSanity *int16
	TempInsanity  *bool
	IndefInsanity *bool
}

type DeleteSanityInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
	MaxMp       *int16
	CurrentMp   *int16
}

type DeleteMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetLuckInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertLuckInput struct {
	UserID       string
	CharacterID  pgtype.UUID
	StartingLuck *int16
	CurrentLuck  *int16
}

type DeleteLuckInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetFinancesInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertFinancesInput struct {
	UserID              string
	CharacterID         pgtype.UUID
	SpendingLimit       *string
	Cash                *string
	Assets              *string
	CreditRatingSkillID pgtype.UUID
}

type DeleteFinancesInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetBackstoryInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertBackstoryInput struct {
	UserID              string
	CharacterID         pgtype.UUID
	PersonalDescription *string
}

type DeleteBackstoryInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetBackstoryItemsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetBackstoryItemInput struct {
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}

type CreateBackstoryItemInput struct {
	Section     string
	Title       string
	Text        string
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateBackstoryItemInput struct {
	Section         string
	Title           string
	Text            string
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}

type DeleteBackstoryItemInput struct {
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}

type GetSkillsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetSkillInput struct {
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}

type CreateSkillInput struct {
	Name        string
	CategoryID  pgtype.UUID
	BaseValue   int16
	Value       int16
	Checked     bool
	Specialized bool
	SpecialtyID pgtype.UUID
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateSkillInput struct {
	Name        string
	CategoryID  pgtype.UUID
	BaseValue   int16
	Value       int16
	Checked     bool
	Specialized bool
	SpecialtyID pgtype.UUID
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}

type DeleteSkillInput struct {
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}

type GetDerivedStatsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertDerivedStatsInput struct {
	UserID      string
	CharacterID pgtype.UUID
	Speed       *int16
	Physique    *int16
	DamageBonus *string
	DodgeValue  *int16
}

type DeleteDerivedStatsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetCharacteristicsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertCharacteristicsInput struct {
	Strength     *int16
	Constitution *int16
	Size         *int16
	Dexterity    *int16
	Appearance   *int16
	Intelligence *int16
	Power        *int16
	Education    *int16
	UserID       string
	CharacterID  pgtype.UUID
}

type DeleteCharacteristicsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetNotesInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetNoteInput struct {
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}

type CreateNoteInput struct {
	Title       string
	Body        string
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateNoteInput struct {
	Title       string
	Body        string
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}

type DeleteNoteInput struct {
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}
