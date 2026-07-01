package finances

import "github.com/jackc/pgx/v5/pgtype"

type FinancesRequest struct {
	SpendingLimit       *string     `json:"spending_limit"`
	Cash                *string     `json:"cash"`
	Assets              *string     `json:"assets"`
	CreditRatingSkillID pgtype.UUID `json:"credit_rating_skill_id"`
}
