package httpAdapter

import (
	"encoding/json"
	"net/http"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) *errors.AppError

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appErr := h(w, r)
	if appErr == nil {
		return
	}
	appErr = errors.NormalizeAppError(appErr)

	middleware.SetAppError(r.Context(), appErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode())
	_ = json.NewEncoder(w).Encode(appErr.Response())
}
