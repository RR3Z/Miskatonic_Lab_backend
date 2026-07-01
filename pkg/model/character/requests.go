package character

import "github.com/jackc/pgx/v5/pgtype"

type CharacterRequest struct {
	Name       string  `json:"name"`
	PlayerName *string `json:"player_name"`
	Occupation *string `json:"occupation"`
	Age        *int16  `json:"age"`
	Sex        *string `json:"sex"`
	Residence  *string `json:"residence"`
	Birthplace *string `json:"birthplace"`
}

type HealthRequest struct {
	MaxHp       *int16 `json:"max_hp"`
	CurrentHp   *int16 `json:"current_hp"`
	MajorWound  *bool  `json:"major_wound"`
	Unconscious *bool  `json:"unconscious"`
	Dying       *bool  `json:"dying"`
	Dead        *bool  `json:"dead"`
}

type SanityRequest struct {
	MaxSanity     *int16 `json:"max_sanity"`
	CurrentSanity *int16 `json:"current_sanity"`
	TempInsanity  *bool  `json:"temp_insanity"`
	IndefInsanity *bool  `json:"indef_insanity"`
}

type MagicRequest struct {
	MaxMp     *int16 `json:"max_mp"`
	CurrentMp *int16 `json:"current_mp"`
}

type LuckRequest struct {
	StartingLuck *int16 `json:"starting_luck"`
	CurrentLuck  *int16 `json:"current_luck"`
}

type FinancesRequest struct {
	SpendingLimit       *string     `json:"spending_limit"`
	Cash                *string     `json:"cash"`
	Assets              *string     `json:"assets"`
	CreditRatingSkillID pgtype.UUID `json:"credit_rating_skill_id"`
}

type BackstoryRequest struct {
	PersonalDescription *string `json:"personal_description"`
}

type BackstoryItemRequest struct {
	Section string `json:"section"`
	Title   string `json:"title"`
	Text    string `json:"text"`
}

type SkillRequest struct {
	Name        string      `json:"name"`
	CategoryID  pgtype.UUID `json:"category_id"`
	BaseValue   int16       `json:"base_value"`
	Value       int16       `json:"value"`
	Checked     bool        `json:"checked"`
	Specialized bool        `json:"specialized"`
	SpecialtyID pgtype.UUID `json:"specialty_id"`
}

type DerivedStatsRequest struct {
	Speed       *int16  `json:"speed"`
	Physique    *int16  `json:"physique"`
	DamageBonus *string `json:"damage_bonus"`
	DodgeValue  *int16  `json:"dodge_value"`
}

type CharacteristicsRequest struct {
	Strength     *int16 `json:"strength"`
	Constitution *int16 `json:"constitution"`
	Size         *int16 `json:"size"`
	Dexterity    *int16 `json:"dexterity"`
	Appearance   *int16 `json:"appearance"`
	Intelligence *int16 `json:"intelligence"`
	Power        *int16 `json:"power"`
	Education    *int16 `json:"education"`
}

type NoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
