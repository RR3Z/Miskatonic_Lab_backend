package financesDTO

import "github.com/jackc/pgx/v5/pgtype"

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
