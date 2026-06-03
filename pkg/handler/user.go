package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

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

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Error(
			"failed to get clerk session claims",
			"component", "user_api",
		)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.Subject
	if strings.TrimSpace(userID) == "" {
		slog.Error(
			"clerk session claims missing subject",
			"component", "user_api",
		)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.services.User.GetUserByID(r.Context(), userID)
	// DB is fine BUT user wasn't find
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error(
			"user not found by clerk id",
			"component", "user_api",
			"clerk_user_id", userID,
			"error", err,
		)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	// Other errors
	if err != nil {
		slog.Error(
			"failed to get user by clerk id",
			"component", "user_api",
			"clerk_user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleUserClerkWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(
			"failed to read clerk user webhook request body",
			"component", "clerk_webhook",
			"error", err,
		)
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	if err := verifyClerkWebhook(payload, r.Header); err != nil {
		slog.Error(
			"failed to verify clerk user webhook signature",
			"component", "clerk_webhook",
			"error", err,
		)
		http.Error(w, "invalid webhook signature", http.StatusUnauthorized)
		return
	}

	var event ClerkWebhookUserEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		slog.Error(
			"failed to parse clerk user webhook payload",
			"component", "clerk_webhook",
			"error", err,
		)
		http.Error(w, "invalid webhook payload", http.StatusBadRequest)
		return
	}

	if event.Type == "user.created" {
		h.createUser(w, r, event.Data)
	} else if event.Type == "user.updated" {
		h.updateUser(w, r, event.Data)
	} else if event.Type == "user.deleted" {
		h.deleteUser(w, r, event.Data.ID)
	} else {
		slog.Error(
			"unexpected clerk user webhook event type",
			"component", "clerk_webhook",
			"event_type", event.Type,
		)
		http.Error(w, "unexpected webhook event type", http.StatusBadRequest)
		return
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request, data ClerkWebhookUserData) {
	input := db.UpsertUserParams{
		ID:        data.ID,
		Username:  parseClerkWebhookUsername(data),
		Email:     parseClerkWebhookEmail(data),
		AvatarUrl: data.ImageURL,
	}

	if err := h.services.User.UpsertUser(r.Context(), input); err != nil {
		slog.Error(
			"failed to create user from clerk webhook",
			"component", "clerk_webhook",
			"clerk_user_id", data.ID,
			"error", err,
		)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request, data ClerkWebhookUserData) {
	input := db.UpsertUserParams{
		ID:        data.ID,
		Username:  parseClerkWebhookUsername(data),
		Email:     parseClerkWebhookEmail(data),
		AvatarUrl: data.ImageURL,
	}

	if err := h.services.User.UpsertUser(r.Context(), input); err != nil {
		slog.Error(
			"failed to update user from clerk webhook",
			"component", "clerk_webhook",
			"clerk_user_id", data.ID,
			"error", err,
		)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, userID string) {
	if strings.TrimSpace(userID) == "" {
		slog.Error(
			"failed to delete user from clerk webhook: missing clerk user id",
			"component", "clerk_webhook",
		)
		http.Error(w, "missing clerk user id", http.StatusBadRequest)
		return
	}

	if err := h.services.User.DeleteUser(r.Context(), userID); err != nil {
		slog.Error(
			"failed to delete user from clerk webhook",
			"component", "clerk_webhook",
			"clerk_user_id", userID,
			"error", err,
		)
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
