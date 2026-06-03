package model

import "github.com/jackc/pgx/v5/pgtype"

type FinancesModel struct {
	ID pgtype.UUID `json:"id"`

	SpendingLimit *string     `json:"spending_limit"`
	Cash          *string     `json:"cash"`
	Assets        *string     `json:"assets"`
	CreditRating  *SkillModel `json:"credit_rating,omitempty"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
