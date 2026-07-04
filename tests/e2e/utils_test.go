package tests

import (
	"crypto/sha1"
	"encoding/hex"
)

func collectIDs(values []e2eIDResponse) []string {
	ids := make([]string, 0, len(values))
	for _, value := range values {
		ids = append(ids, value.ID)
	}
	return ids
}

func collectSelectedCharacterIDs(values []e2eSelectedCharacterResponse) []string {
	ids := make([]string, 0, len(values))
	for _, value := range values {
		ids = append(ids, value.Character.ID)
	}
	return ids
}

func e2eHash(value string) string {
	sum := sha1.Sum([]byte(value))
	return hex.EncodeToString(sum[:])[:12]
}
