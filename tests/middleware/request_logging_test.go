package tests

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/stretchr/testify/require"
)

func TestSetAppErrorNoOpsWithoutLogState(t *testing.T) {
	require.NotPanics(t, func() {
		middleware.SetAppError(context.Background(), &myErrors.AppError{Code: "test"})
	})
}

func TestRequestLoggingMiddlewareLogsStatusLevelAndErrorCode(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		appErr    *myErrors.AppError
		wantLevel string
		wantText  string
		wantCode  string
	}{
		{
			name:      "success default status",
			status:    0,
			wantLevel: `"level":"INFO"`,
			wantText:  `"msg":"request completed"`,
		},
		{
			name:   "client error",
			status: http.StatusBadRequest,
			appErr: &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Code:    "character.invalid_input",
				Message: "invalid request body",
			},
			wantLevel: `"level":"WARN"`,
			wantText:  `"msg":"invalid request body"`,
			wantCode:  `"error_code":"character.invalid_input"`,
		},
		{
			name:   "server error",
			status: http.StatusInternalServerError,
			appErr: &myErrors.AppError{
				Status:  http.StatusInternalServerError,
				Code:    "common.internal_error",
				Message: "failed to get character",
			},
			wantLevel: `"level":"ERROR"`,
			wantText:  `"msg":"failed to get character"`,
			wantCode:  `"error_code":"common.internal_error"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var logs bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			handler := middleware.RequestLoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.appErr != nil {
					middleware.SetAppError(r.Context(), tc.appErr)
				}
				if tc.status != 0 {
					w.WriteHeader(tc.status)
				}
			}))

			handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", nil))

			output := logs.String()
			require.Contains(t, output, tc.wantLevel)
			require.Contains(t, output, tc.wantText)
			require.Contains(t, output, `"method":"GET"`)
			require.Contains(t, output, `"path":"/test"`)
			if tc.wantCode != "" {
				require.Contains(t, output, tc.wantCode)
			}
		})
	}
}

func TestRequestLoggingMiddlewareForwardsHijacker(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() { _ = serverConn.Close() })
	t.Cleanup(func() { _ = clientConn.Close() })

	writer := &hijackableResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		conn:           serverConn,
	}
	handler := middleware.RequestLoggingMiddleware(slog.New(slog.NewTextHandler(io.Discard, nil)))(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hijacker, ok := w.(http.Hijacker)
		require.True(t, ok)

		conn, _, err := hijacker.Hijack()
		require.NoError(t, err)
		require.Same(t, serverConn, conn)
	}))

	handler.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/ws", nil))
	require.True(t, writer.hijacked)
}

type hijackableResponseWriter struct {
	http.ResponseWriter
	conn     net.Conn
	hijacked bool
}

func (w *hijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	return w.conn, bufio.NewReadWriter(bufio.NewReader(w.conn), bufio.NewWriter(w.conn)), nil
}
