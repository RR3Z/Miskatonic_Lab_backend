package errors

import "errors"

var ErrCurrentHealthExceedsMax = errors.New("current_hp cannot exceed max_hp")
