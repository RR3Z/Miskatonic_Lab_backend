package portrait

import "errors"

var (
	ErrPortraitRequired   = errors.New("portrait is required")
	ErrPortraitTooLarge   = errors.New("portrait is too large")
	ErrUnsupportedImage   = errors.New("unsupported portrait image")
	ErrInvalidImage       = errors.New("invalid portrait image")
	ErrInvalidPortraitKey = errors.New("invalid portrait key")
)
