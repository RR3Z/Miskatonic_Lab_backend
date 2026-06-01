package MiskatonicLab

import "time"

type User struct {
	Id int `json:"-"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
