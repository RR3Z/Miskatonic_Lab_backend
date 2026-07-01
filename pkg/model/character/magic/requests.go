package magicDTO

type MagicRequest struct {
	MaxMp     *int16 `json:"max_mp"`
	CurrentMp *int16 `json:"current_mp"`
}
