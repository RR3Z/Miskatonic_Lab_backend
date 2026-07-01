package characters

type CharacterRequest struct {
	Name       string  `json:"name"`
	PlayerName *string `json:"player_name"`
	Occupation *string `json:"occupation"`
	Age        *int16  `json:"age"`
	Sex        *string `json:"sex"`
	Residence  *string `json:"residence"`
	Birthplace *string `json:"birthplace"`
}
