package room

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid room input")
	ErrRoomNotFound      = errors.New("room not found")
	ErrRoomFull          = errors.New("room is full")
	ErrAlreadyMember     = errors.New("already a member of this room")
	ErrNotMember         = errors.New("not a member of this room")
	ErrNotOwner          = errors.New("only the room owner can perform this action")
	ErrCannotKickOwner   = errors.New("cannot kick the room owner")
	ErrCharacterNotOwned = errors.New("character does not belong to you")
)
