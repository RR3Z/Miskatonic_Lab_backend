package room

// Default Settings for Room
const (
	DefaultMaxPlayers int32 = 7
)

// Players Role
const (
	RolePlayer = "player"
	RoleGM     = "gm"
)

func IsValidRole(role string) bool {
	return role == RolePlayer || role == RoleGM
}
