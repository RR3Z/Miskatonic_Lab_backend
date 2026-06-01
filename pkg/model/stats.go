package MiskatonicLab

import "time"

type Characteristics struct {
	Id string `json:"-"`

	Strength     *int `json:"str,omitempty"` // Strength, Сила
	Constitution *int `json:"con,omitempty"` // Constitution, Выносливость
	Size         *int `json:"siz,omitempty"` // Size, Телосложение
	Dexterity    *int `json:"dex,omitempty"` // Dexterity, Ловкость
	Appearance   *int `json:"app,omitempty"` // Appearance, Наружность
	Intelligence *int `json:"int,omitempty"` // Intelligence, Интеллект
	Power        *int `json:"pow,omitempty"` // Power, Мощь
	Education    *int `json:"edu,omitempty"` // Education, Образование

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DerivedStats struct {
	Id string `json:"-"`

	Speed       *int `json:"speed,omitempty"`
	Physique    *int `json:"physique,omitempty"` // Physique, Комплекция
	DamageBonus *int `json:"damageBonus,omitempty"`
	DodgeValue  *int `json:"dodgeValue,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
