package errors

import "errors"

var ErrCurrentSanityExceedsMax = errors.New("current_sanity cannot exceed max_sanity")
