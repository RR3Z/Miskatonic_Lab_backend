package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	handlerErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/user/errors"
	userHandlerHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/user/helpers"
	userDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/clerk/clerk-sdk-go/v2"
)

func (h *UserHandler) getMe(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		return &myErrors.AppError{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Err:     errors.New("failed to get clerk session claims"),
		}
	}

	userID := claims.Subject
	if strings.TrimSpace(userID) == "" {
		return &myErrors.AppError{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Err:     errors.New("clerk session claims missing subject"),
		}
	}

	user, err := h.service.GetUserByID(r.Context(), userDTO.GetUserInput{ID: userID})
	if err != nil {
		if errors.Is(err, userErrors.ErrUserNotFound) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}
		}
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get user",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, user)
	return nil
}

func (h *UserHandler) handleClerkWebhook(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return handlerErrors.InvalidRequestBodyError(err)
	}

	if err := userHandlerHelpers.VerifyClerkWebhook(payload, r.Header); err != nil {
		return handlerErrors.InvalidWebhookSignatureError(err)
	}

	var event userDTO.ClerkWebhookUserEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return handlerErrors.InvalidWebhookPayloadError(err)
	}

	switch event.Type {
	case "user.created":
		return h.createUserFromWebhook(w, r, event.Data)
	case "user.updated":
		return h.updateUserFromWebhook(w, r, event.Data)
	case "user.deleted":
		return h.deleteUserFromWebhook(w, r, event.Data.ID)
	default:
		return handlerErrors.UnexpectedWebhookEventError(event.Type)
	}
}

func (h *UserHandler) createUserFromWebhook(w http.ResponseWriter, r *http.Request, data userDTO.ClerkWebhookUserData) *myErrors.AppError {
	input := userDTO.ToUpsertUserInput(data)

	if err := h.service.UpsertUser(r.Context(), input); err != nil {
		return handlerErrors.MapServiceError(err, "failed to create user")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *UserHandler) updateUserFromWebhook(w http.ResponseWriter, r *http.Request, data userDTO.ClerkWebhookUserData) *myErrors.AppError {
	input := userDTO.ToUpsertUserInput(data)

	if err := h.service.UpsertUser(r.Context(), input); err != nil {
		return handlerErrors.MapServiceError(err, "failed to update user")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *UserHandler) deleteUserFromWebhook(w http.ResponseWriter, r *http.Request, userID string) *myErrors.AppError {
	if err := h.service.DeleteUser(r.Context(), userDTO.DeleteUserInput{ID: userID}); err != nil {
		return handlerErrors.MapServiceError(err, "failed to delete user")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
