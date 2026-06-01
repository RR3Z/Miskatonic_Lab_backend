package MiskatonicLab

import "time"

type Skill struct {
	Id string `json:"-"`

	Name     string        `json:"name"`
	Category SkillCategory `json:"category"`

	BaseValue int  `json:"baseValue"`
	Value     int  `json:"value"`
	Checked   bool `json:"checked"`

	Specialized bool       `json:"specialized,omitempty"`
	Specialty   *Specialty `json:"specialty,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SkillCategory struct {
	Id   string `json:"-"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Specialty struct {
	Id string `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`

	BaseValue int `json:"baseValue"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
