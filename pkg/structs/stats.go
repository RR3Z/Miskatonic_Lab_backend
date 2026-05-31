package MiskatonicLab

import "time"

type Characteristics struct {
	STR int `json:"str"` // Strength, Сила
	CON int `json:"con"` // Constitution, Выносливость
	SIZ int `json:"siz"` // Size, Телосложение
	DEX int `json:"dex"` // Dexterity, Ловкость
	APP int `json:"app"` // Appearance, Наружность
	INT int `json:"int"` // Intelligence, Интеллект
	POW int `json:"pow"` // Power, Мощь
	EDU int `json:"edu"` // Education, Образование

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DerivedStats struct {
	Speed       int    `json:"speed"`
	Build       int    `json:"build"` // Physique, Комплекция
	DamageBonus string `json:"damageBonus"`
	DodgeBase   int    `json:"dodgeBase"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
