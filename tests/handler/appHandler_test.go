package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/stretchr/testify/require"
)

func TestAppHandlerWritesJSONErrorResponse(t *testing.T) {
	subject := handler.AppHandler(func(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Code:    "room.not_found",
			Message: "room not found",
		}
	})

	recorder := httptest.NewRecorder()
	subject.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	require.JSONEq(t, `{"code":"room.not_found","message":"room not found"}`, recorder.Body.String())
}

func TestAppHandlerDefaultsErrorCodeWhenMissing(t *testing.T) {
	subject := handler.AppHandler(func(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		}
	})

	recorder := httptest.NewRecorder()
	subject.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/test", nil))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.JSONEq(t, `{"code":"common.invalid_request","message":"invalid request body"}`, recorder.Body.String())
}

func TestAppHandlerStopsAfterError(t *testing.T) {
	continued := false
	subject := handler.AppHandler(func(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
		if true {
			return &myErrors.AppError{
				Status:  http.StatusBadRequest,
				Code:    "common.invalid_request",
				Message: "invalid request",
			}
		}

		continued = true
		return nil
	})

	recorder := httptest.NewRecorder()
	subject.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.False(t, continued)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
