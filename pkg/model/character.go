package MiskatonicLab

import "time"

type Character struct {
	Id string `json:"-"`

	Name       string  `json:"name"`
	PlayerName *string `json:"playerName,omitempty"`
	Occupation *string `json:"occupation,omitempty"`
	Age        *int    `json:"age,omitempty"`
	Sex        *string `json:"sex,omitempty"`
	Residence  *string `json:"residence,omitempty"`
	Birthplace *string `json:"birthplace,omitempty"`

	Characteristics Characteristics `json:"characteristics"`
	DerivedStats    DerivedStats    `json:"derivedStats"`
	Skills          []Skill         `json:"skills"`
	Finances        Finances        `json:"finances"`

	HealthState HealthState `json:"healthState"`
	SanityState SanityState `json:"sanityState"`
	MagicState  MagicState  `json:"magicState"`
	LuckState   LuckState   `json:"luckState"`

	Backstory Backstory `json:"backstory"`

	Notes []Note `json:"notes"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Note struct {
	Id string `json:"-"`

	Title string `json:"title"`
	Body  string `json:"body"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Finances struct {
	Id string `json:"-"`

	SpendingLimit *string `json:"spendingLevel,omitempty"`
	Cash          *string `json:"cash,omitempty"`
	Assets        *string `json:"assets,omitempty"`

	CreditRating *Skill `json:"creditRating,omitempty"` // Credit Rating, Средства

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
