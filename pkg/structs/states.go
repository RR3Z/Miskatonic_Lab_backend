package MiskatonicLab

import "time"

type HealthState struct {
	CurrentHitPoints int `json:"currentHitPoints"`
	MaxHitPoints     int `json:"maxHitPoints"`

	MajorWound bool `json:"majorWound"`

	Unconscious bool `json:"unconscious"`
	Dying       bool `json:"dying"`
	Dead        bool `json:"dead"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SanityState struct {
	CurrentSanity int  `json:"currentSanity"`
	MaxSanity     int  `json:"maxSanity"`
	TempInsanity  bool `json:"tempInsanity"`
	IndefInsanity bool `json:"indefInsanity"` // Indefinite Insanity, Продолжительное безумие

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MagicState struct {
	CurrentMagicPoints int `json:"currentMagicPoints"`
	MaxMagicPoints     int `json:"maxMagicPoints"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LuckState struct {
	CurrentLuck  int `json:"currentLuck"`
	StartingLuck int `json:"startingLuck"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
