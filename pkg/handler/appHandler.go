package handler

import (
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

	middleware.SetAppError(r.Context(), appErr)

	http.Error(w, appErr.Message, appErr.Status)
}
