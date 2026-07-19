package roomDTO

type RoomMutationResult[T any] struct {
	Value  T
	Events []RoomEventModel
}
