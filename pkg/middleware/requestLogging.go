package middleware

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
)

type RequestLogState struct {
	AppErr *errors.AppError
}

type RequestLogStateKey struct{}

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("wrapped response writer does not implement http.Hijacker")
	}

	return hijacker.Hijack()
}

func SetAppError(ctx context.Context, appErr *errors.AppError) {
	state, ok := ctx.Value(RequestLogStateKey{}).(*RequestLogState)
	if !ok || state == nil {
		return
	}

	state.AppErr = appErr
}

func RequestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			state := &RequestLogState{}
			ctx := context.WithValue(r.Context(), RequestLogStateKey{}, state)
			r = r.WithContext(ctx)

			rw := &ResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", time.Since(startedAt).String(),
			}

			if state.AppErr != nil && state.AppErr.Err != nil {
				attrs = append(attrs, "error", state.AppErr.Err)
			}
			if state.AppErr != nil && state.AppErr.Response().Code != "" {
				attrs = append(attrs, "error_code", state.AppErr.Response().Code)
			}

			message := "request completed"
			if rw.status >= 400 {
				message = "request failed"
			}
			if state.AppErr != nil && state.AppErr.Message != "" {
				message = state.AppErr.Message
			}

			switch {
			case rw.status >= 500:
				logger.ErrorContext(r.Context(), message, attrs...)
			case rw.status >= 400:
				logger.WarnContext(r.Context(), message, attrs...)
			default:
				logger.InfoContext(r.Context(), message, attrs...)
			}
		})
	}
}
