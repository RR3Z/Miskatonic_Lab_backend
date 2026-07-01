package luckDTO

type LuckRequest struct {
	StartingLuck *int16 `json:"starting_luck"`
	CurrentLuck  *int16 `json:"current_luck"`
}
