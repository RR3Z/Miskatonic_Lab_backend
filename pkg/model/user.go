package MiskatonicLab

import "time"

type User struct {
	Id          string `json:"-"`
	ClerkUserId string `json:"-"`

	Username  string  `json:"username"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatarUrl,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
