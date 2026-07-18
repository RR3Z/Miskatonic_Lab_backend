package financesDTO

type FinancesRequest struct {
	SpendingLimit *string `json:"spending_limit"`
	Cash          *string `json:"cash"`
	Assets        *string `json:"assets"`
}
