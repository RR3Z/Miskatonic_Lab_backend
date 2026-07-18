package financesDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type FinancesModel struct {
	ID pgtype.UUID `json:"id"`

	SpendingLimit *string `json:"spending_limit"`
	Cash          *string `json:"cash"`
	Assets        *string `json:"assets"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func ToFinancesModel(f db.Finance) FinancesModel {
	return FinancesModel{
		ID:            f.ID,
		SpendingLimit: f.SpendingLimit,
		Cash:          f.Cash,
		Assets:        f.Assets,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}
