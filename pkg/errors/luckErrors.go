package errors

import "errors"

var ErrCurrentLuckExceedsStarting = errors.New("current_luck cannot exceed starting_luck")
