package room

// Default Settings for Room
const (
	DEFAULT_MAX_PLAYERS int32 = 7
)

const (
	DEFAULT_ROOM_EVENTS_LIMIT int32 = 100
	MAX_ROOM_EVENTS_LIMIT     int32 = 200
	MAX_CHAT_MESSAGE_LENGTH   int   = 2000
)

// Players Role
const (
	ROLE_PLAYER = "player"
	ROLE_GM     = "gm"
)

func IsValidRole(role string) bool {
	return role == ROLE_PLAYER || role == ROLE_GM
}
