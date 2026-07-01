package derivedStatsDTO

import "github.com/jackc/pgx/v5/pgtype"

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
