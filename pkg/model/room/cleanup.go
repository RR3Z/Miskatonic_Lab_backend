package roomDTO

type CleanupRoomsResult struct {
	InactiveDeleted int `json:"inactive_deleted"`
	InvalidDeleted  int `json:"invalid_deleted"`
}
