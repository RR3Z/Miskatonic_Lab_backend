package characteristics

type CharacteristicsRequest struct {
	Strength     *int16 `json:"strength"`
	Constitution *int16 `json:"constitution"`
	Size         *int16 `json:"size"`
	Dexterity    *int16 `json:"dexterity"`
	Appearance   *int16 `json:"appearance"`
	Intelligence *int16 `json:"intelligence"`
	Power        *int16 `json:"power"`
	Education    *int16 `json:"education"`
}
