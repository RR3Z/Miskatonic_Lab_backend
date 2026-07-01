package sanityDTO

type SanityRequest struct {
	MaxSanity     *int16 `json:"max_sanity"`
	CurrentSanity *int16 `json:"current_sanity"`
	TempInsanity  *bool  `json:"temp_insanity"`
	IndefInsanity *bool  `json:"indef_insanity"`
}
