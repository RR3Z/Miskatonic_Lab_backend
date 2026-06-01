package MiskatonicLab

import "time"

type Backstory struct {
	Id string `json:"-"`

	PersonalDescription *string `json:"personalDescription,omitempty"`

	InjuriesScars        []BackstoryItem `json:"injuriesScars"`
	PhobiasManias        []BackstoryItem `json:"phobiasManias"`
	ArcaneTomesSpells    []BackstoryItem `json:"arcaneTomesSpells"`
	Encounters           []BackstoryItem `json:"encounters"`
	IdeologyBeliefs      []BackstoryItem `json:"ideologyBeliefs"`
	SignificantPeople    []BackstoryItem `json:"significantPeople"`
	MeaningfulLocations  []BackstoryItem `json:"meaningfulLocations"`
	TreasuredPossessions []BackstoryItem `json:"treasuredPossessions"`
	Traits               []BackstoryItem `json:"traits"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BackstoryItem struct {
	Id string `json:"-"`

	Title string `json:"title"`
	Text  string `json:"text"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
