package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
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

func (s *e2eSubject) deleteCharacter(t *testing.T, characterID string) {
	t.Helper()
	s.doJSONAllow(t, http.MethodDelete, "/api/characters/"+url.PathEscape(characterID)+"/", nil, []int{http.StatusNoContent, http.StatusNotFound}, nil)
}

func (s *e2eSubject) createRoom(t *testing.T, password string) e2eRoomResponse {
	t.Helper()

	var room e2eRoomResponse
	s.doJSON(
		t,
		http.MethodPost,
		"/api/rooms/",
		map[string]any{"max_players": 4, "password": password},
		http.StatusCreated,
		&room,
	)
	require.NotEmpty(t, room.ID)
	return room
}

func (s *e2eSubject) deleteRoom(t *testing.T, roomID string) {
	t.Helper()
	s.doJSONAllow(t, http.MethodDelete, "/api/rooms/"+url.PathEscape(roomID)+"/", nil, []int{http.StatusNoContent, http.StatusNotFound}, nil)
}

func (s *e2eSubject) joinRoom(t *testing.T, roomID string, password string) {
	t.Helper()
	s.doJSON(
		t,
		http.MethodPost,
		"/api/rooms/"+url.PathEscape(roomID)+"/join",
		map[string]any{"password": password},
		http.StatusOK,
		nil,
	)
}

func (s *e2eSubject) selectCharacter(t *testing.T, roomID string, characterID string) {
	t.Helper()
	s.doJSON(
		t,
		http.MethodPut,
		"/api/rooms/"+url.PathEscape(roomID)+"/character",
		map[string]any{"character_id": characterID},
		http.StatusOK,
		nil,
	)
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

func (s *e2eSubject) doMultipartFile(t *testing.T, method string, path string, field string, filename string, content []byte, expectedStatus int, target any) {
	t.Helper()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	var part io.Writer
	var err error
	if filename == "" {
		part, err = writer.CreateFormField(field)
	} else {
		part, err = writer.CreateFormFile(field, filename)
	}
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(method, s.baseURL+path, body)
	require.NoError(t, err)
	req.Header.Set("Authorization", s.authorization(t))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := s.client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatus, res.StatusCode, "response body: %s", string(responseBody))
	if target != nil && len(responseBody) > 0 {
		require.NoError(t, json.Unmarshal(responseBody, target), "response body: %s", string(responseBody))
	}
}

func (s *e2eSubject) doPublicPortraitRequest(t *testing.T, method string, portraitURL string, expectedStatus int) ([]byte, http.Header) {
	t.Helper()

	req, err := http.NewRequest(method, portraitURL, nil)
	require.NoError(t, err)
	res, err := s.client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatus, res.StatusCode, "response body: %s", string(body))
	return body, res.Header
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
	req.Header.Set("Authorization", s.authorization(t))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func (s *e2eSubject) wsURL(t *testing.T, path string) string {
	t.Helper()

	value, err := url.Parse(s.baseURL + path)
	require.NoError(t, err)
	switch value.Scheme {
	case "http":
		value.Scheme = "ws"
	case "https":
		value.Scheme = "wss"
	default:
		t.Fatalf("unsupported E2E_BASE_URL scheme %q", value.Scheme)
	}
	return value.String()
}

func (s *e2eSubject) authorization(t testing.TB) string {
	t.Helper()
	require.NotNil(t, suiteE2EClerkFixture, "E2E Clerk fixture is not initialized")
	return suiteE2EClerkFixture.authorization(t, s.identity)
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
