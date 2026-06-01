package utils

import "strings"

func ParseAllowedOrigins(value string) []string {
	origins := make([]string, 0)

	for _, origin := range strings.Split(value, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
