package errors

import "errors"

var ErrCurrentMagicExceedsMax = errors.New("current_mp cannot exceed max_mp")
