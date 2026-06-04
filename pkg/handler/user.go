package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5"
	svix "github.com/svix/svix-webhooks/go"
)

var errMissingClerkWebhookSigningSecret = errors.New("CLERK_WEBHOOK_SIGNING_SECRET is not set")

type ClerkWebhookUserEvent struct {
	Type string               `json:"type"`
	Data ClerkWebhookUserData `json:"data"`
}

type ClerkWebhookUserData struct {
	ID                    string                  `json:"id"`
	Username              *string                 `json:"username"`
	ImageURL              *string                 `json:"image_url"`
	PrimaryEmailAddressID *string                 `json:"primary_email_address_id"`
	EmailAddresses        []ClerkWebhookUserEmail `json:"email_addresses"`
}

type ClerkWebhookUserEmail struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
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

	user, err := h.services.User.GetUserByID(r.Context(), userID)
	// DB is fine BUT user wasn't find
	if errors.Is(err, pgx.ErrNoRows) {
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Message: "user not found",
			Err:     err,
		}
	}
	// Other errors
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get user",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, user)
	return nil
}

func (h *Handler) handleUserClerkWebhook(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "failed to read request body",
			Err:     err,
		}
	}

	if err := verifyClerkWebhook(payload, r.Header); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusUnauthorized,
			Message: "invalid webhook signature",
			Err:     err,
		}
	}

	var event ClerkWebhookUserEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid webhook payload",
			Err:     err,
		}
	}

	switch event.Type {
	case "user.created":
		return h.createUser(w, r, event.Data)
	case "user.updated":
		return h.updateUser(w, r, event.Data)
	case "user.deleted":
		return h.deleteUser(w, r, event.Data.ID)
	default:
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "unexpected webhook event type",
			Err:     errors.New("unexpected clerk user webhook event type: " + event.Type),
		}
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request, data ClerkWebhookUserData) *myErrors.AppError {
	input := db.UpsertUserParams{
		ID:        data.ID,
		Username:  parseClerkWebhookUsername(data),
		Email:     parseClerkWebhookEmail(data),
		AvatarUrl: data.ImageURL,
	}

	if err := h.services.User.UpsertUser(r.Context(), input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create user",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request, data ClerkWebhookUserData) *myErrors.AppError {
	input := db.UpsertUserParams{
		ID:        data.ID,
		Username:  parseClerkWebhookUsername(data),
		Email:     parseClerkWebhookEmail(data),
		AvatarUrl: data.ImageURL,
	}

	if err := h.services.User.UpsertUser(r.Context(), input); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update user",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, userID string) *myErrors.AppError {
	if strings.TrimSpace(userID) == "" {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "missing clerk user id",
			Err:     errors.New("failed to delete user from clerk webhook: missing clerk user id"),
		}
	}

	if err := h.services.User.DeleteUser(r.Context(), userID); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete user",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func verifyClerkWebhook(payload []byte, headers http.Header) error {
	signingSecret := os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET")
	if signingSecret == "" {
		return errMissingClerkWebhookSigningSecret
	}

	webhook, err := svix.NewWebhook(signingSecret)
	if err != nil {
		return err
	}

	return webhook.Verify(payload, headers)
}

func parseClerkWebhookUsername(userData ClerkWebhookUserData) string {
	if userData.Username != nil && strings.TrimSpace(*userData.Username) != "" {
		return strings.TrimSpace(*userData.Username)
	}

	email := parseClerkWebhookEmail(userData)
	if email != "" {
		return strings.Split(email, "@")[0]
	}

	return userData.ID
}

func parseClerkWebhookEmail(userData ClerkWebhookUserData) string {
	if userData.PrimaryEmailAddressID != nil {
		for _, email := range userData.EmailAddresses {
			if email.ID == *userData.PrimaryEmailAddressID {
				return email.EmailAddress
			}
		}
	}

	if len(userData.EmailAddresses) > 0 {
		return userData.EmailAddresses[0].EmailAddress
	}

	return userData.ID + "@users.local"
}
