package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func (s *e2eSubject) createCharacter(t *testing.T, name string) string {
	t.Helper()

	var character e2eIDResponse
	s.doJSON(
		t,
		http.MethodPost,
		"/api/characters/",
		map[string]any{"name": name, "age": 37},
		http.StatusCreated,
		&character,
	)
	require.NotEmpty(t, character.ID)
	return character.ID
}

func (s *e2eSubject) doJSON(t *testing.T, method string, path string, body any, expectedStatus int, target any) {
	t.Helper()
	s.doJSONAllow(t, method, path, body, []int{expectedStatus}, target)
}

func (s *e2eSubject) doJSONAllow(t *testing.T, method string, path string, body any, expectedStatuses []int, target any) {
	t.Helper()

	req := s.newRequest(t, method, path, body)
	res, err := s.client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Containsf(t, expectedStatuses, res.StatusCode, "response body: %s", string(responseBody))

	if target != nil && len(responseBody) > 0 {
		require.NoError(t, json.Unmarshal(responseBody, target), "response body: %s", string(responseBody))
	}
}

func (s *e2eSubject) newRequest(t *testing.T, method string, path string, body any) *http.Request {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, s.baseURL+path, reader)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+s.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func (s *e2eSubject) waitForRoomEvents(t *testing.T, roomID string, eventType string) []e2eRoomEventResponse {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	eventsURL := "/api/rooms/" + url.PathEscape(roomID) + "/events?limit=20"
	for time.Now().Before(deadline) {
		var events []e2eRoomEventResponse
		s.doJSON(t, http.MethodGet, eventsURL, nil, http.StatusOK, &events)
		for _, event := range events {
			if event.Type == eventType {
				return events
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for room event type %q", eventType)
	return nil
}
