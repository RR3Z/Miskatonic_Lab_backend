package model

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type FinancesModel struct {
	ID pgtype.UUID `json:"id"`

	SpendingLimit *string     `json:"spending_limit"`
	Cash          *string     `json:"cash"`
	Assets        *string     `json:"assets"`
	CreditRating  *SkillModel `json:"credit_rating,omitempty"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func ToFinancesModel(f db.Finance, creditRating *SkillModel) FinancesModel {
	return FinancesModel{
		ID:            f.ID,
		SpendingLimit: f.SpendingLimit,
		Cash:          f.Cash,
		Assets:        f.Assets,
		CreditRating:  creditRating,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}
