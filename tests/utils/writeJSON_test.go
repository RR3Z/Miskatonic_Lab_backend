package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestWriteJSONWritesStatusContentTypeAndJSONBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	value := map[string]any{
		"id":   "user_1",
		"name": "Roger",
	}

	utils.WriteJSON(recorder, http.StatusCreated, value)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var body map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"id":   "user_1",
		"name": "Roger",
	}, body)
}

func TestWriteJSONWritesNullForNilValue(t *testing.T) {
	recorder := httptest.NewRecorder()

	utils.WriteJSON(recorder, http.StatusOK, nil)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	require.JSONEq(t, "null", recorder.Body.String())
}

func TestWriteJSONOverwritesExistingContentType(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "text/plain")

	utils.WriteJSON(recorder, http.StatusOK, map[string]string{"ok": "true"})

	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestWriteJSONKeepsStatusWhenHeaderWasAlreadyWritten(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(http.StatusAccepted)

	utils.WriteJSON(recorder, http.StatusCreated, map[string]string{"ok": "true"})

	require.Equal(t, http.StatusAccepted, recorder.Code)
	require.JSONEq(t, `{"ok":"true"}`, recorder.Body.String())
}

func TestWriteJSONDoesNotPanicWhenValueCannotBeEncoded(t *testing.T) {
	recorder := httptest.NewRecorder()

	require.NotPanics(t, func() {
		utils.WriteJSON(recorder, http.StatusOK, map[string]any{
			"invalid": func() {},
		})
	})

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	require.Empty(t, recorder.Body.String())
}

func TestWriteJSONDoesNotPanicWhenWriterReturnsError(t *testing.T) {
	writer := &errorResponseWriter{header: http.Header{}}

	require.NotPanics(t, func() {
		utils.WriteJSON(writer, http.StatusCreated, map[string]string{"ok": "true"})
	})

	require.Equal(t, http.StatusCreated, writer.statusCode)
	require.Equal(t, "application/json", writer.header.Get("Content-Type"))
	require.True(t, writer.writeCalled)
}

type errorResponseWriter struct {
	header      http.Header
	statusCode  int
	writeCalled bool
}

func (w *errorResponseWriter) Header() http.Header {
	return w.header
}

func (w *errorResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *errorResponseWriter) Write([]byte) (int, error) {
	w.writeCalled = true
	return 0, errors.New("write failed")
}
