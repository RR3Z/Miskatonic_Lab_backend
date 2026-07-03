package room

import "strings"

func validateMaxPlayers(maxPlayers int32) error {
	if maxPlayers < 1 {
		return ErrInvalidInput
	}
	return nil
}

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrInvalidPassword
	}
	return nil
}

func validateChatMessage(text string) error {
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" || len(trimmedText) > MAX_CHAT_MESSAGE_LENGTH {
		return ErrInvalidInput
	}
	return nil
}

func normalizeRoomEventsLimit(limit int32) int32 {
	if limit <= 0 {
		return DEFAULT_ROOM_EVENTS_LIMIT
	}
	if limit > MAX_ROOM_EVENTS_LIMIT {
		return MAX_ROOM_EVENTS_LIMIT
	}
	return limit
}
