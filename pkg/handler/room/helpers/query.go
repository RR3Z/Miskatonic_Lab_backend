package roomHelpers

import (
	"net/http"
	"strconv"
	"strings"
)

func OptionalInt32Query(r *http.Request, key string) (*int32, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return nil, err
	}

	result := int32(parsed)
	return &result, nil
}

func OptionalInt64Query(r *http.Request, key string) (*int64, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
