package MiskatonicLab

import "time"

type Backstory struct {
	Id string `json:"-"`

	PersonalDescription string `json:"personalDescription,omitempty"`

	IdeologyBeliefs      []BackstoryItem `json:"ideologyBeliefs,omitempty"`
	SignificantPeople    []BackstoryItem `json:"significantPeople,omitempty"`
	MeaningfulLocations  []BackstoryItem `json:"meaningfulLocations,omitempty"`
	TreasuredPossessions []BackstoryItem `json:"treasuredPossessions,omitempty"`
	Traits               []BackstoryItem `json:"traits,omitempty"`

	InjuriesScars     string `json:"injuriesScars,omitempty"`
	PhobiasManias     string `json:"phobiasManias,omitempty"`
	ArcaneTomesSpells string `json:"arcaneTomesSpells,omitempty"`
	Encounters        string `json:"encounters,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BackstoryItem struct {
	Id string `json:"-"`

	Title bool   `json:"title"`
	Text  string `json:"text"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
