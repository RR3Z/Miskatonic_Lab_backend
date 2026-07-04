package tests

import (
	"bytes"
	"context"
	"log/slog"
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
