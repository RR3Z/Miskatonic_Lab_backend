package MiskatonicLab

import "time"

type HealthState struct {
	Id string `json:"-"`

	MaxHP     int `json:"maxHP"`
	CurrentHP int `json:"currentHP"`

	MajorWound bool `json:"majorWound"`

	Unconscious bool `json:"unconscious"`
	Dying       bool `json:"dying"`
	Dead        bool `json:"dead"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SanityState struct {
	Id string `json:"-"`

	MaxSanity     int `json:"maxSanity"`
	CurrentSanity int `json:"currentSanity"`

	TempInsanity  bool `json:"tempInsanity"`  // Temporary Insanity, Временное безумие
	IndefInsanity bool `json:"indefInsanity"` // Indefinite Insanity, Продолжительное безумие

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MagicState struct {
	Id string `json:"-"`

	MaxMP     int `json:"maxMP"`
	CurrentMP int `json:"currentMP"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LuckState struct {
	Id string `json:"-"`

	StartingLuck int `json:"startingLuck"`
	CurrentLuck  int `json:"currentLuck"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
