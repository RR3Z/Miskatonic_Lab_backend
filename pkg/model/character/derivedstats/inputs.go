package derivedStatsDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetDerivedStatsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
