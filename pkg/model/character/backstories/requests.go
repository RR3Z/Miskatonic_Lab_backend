package backstoriesDTO

type BackstoryRequest struct {
	PersonalDescription *string `json:"personal_description"`
}

type BackstoryItemRequest struct {
	Section string `json:"section"`
	Title   string `json:"title"`
	Text    string `json:"text"`
}
