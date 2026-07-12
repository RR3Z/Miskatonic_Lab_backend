package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/require"
)

const (
	testIssuer = "https://test-instance.clerk.accounts.dev"
	testParty  = "https://app.example.com"
)

func TestClerkAuthMiddlewareAcceptsValidJWTAndSetsClaims(t *testing.T) {
	now := time.Date(2026, 7, 12, 20, 0, 0, 0, time.UTC)
	token, publicKey := generateClerkToken(t, "valid-kid", now, nil)
	client := newTestJWKSClient(t, []jose.JSONWebKey{newTestJWK("valid-kid", publicKey)})

	handler := middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
		JWKSClient:        client,
		AuthorizedParties: []string{testParty},
		Clock:             clerktest.NewClockAt(now),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		require.True(t, ok)
		require.Equal(t, "user_test", claims.Subject)
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := serveWithToken(handler, token)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestClerkAuthMiddlewareKeepsMissingAndMalformedTokenForbidden(t *testing.T) {
	logs := new(bytes.Buffer)
	handler := middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
		Logger: slog.New(slog.NewJSONHandler(logs, nil)),
	})(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	missing := httptest.NewRecorder()
	handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/me", nil))
	require.Equal(t, http.StatusForbidden, missing.Code)

	malformed := serveWithToken(handler, "not-a-jwt")
	require.Equal(t, http.StatusForbidden, malformed.Code)
	require.Equal(t, 2, strings.Count(logs.String(), `"auth_failure":"missing_or_malformed_token"`))
	require.NotContains(t, logs.String(), "not-a-jwt")
}

func TestClerkAuthMiddlewareAllowsConfiguredClockSkewLeeway(t *testing.T) {
	now := time.Date(2026, 7, 12, 20, 0, 0, 0, time.UTC)
	token, publicKey := generateClerkToken(t, "leeway-kid", now, map[string]any{
		"nbf": now.Add(5 * time.Second).Unix(),
	})
	client := newTestJWKSClient(t, []jose.JSONWebKey{newTestJWK("leeway-kid", publicKey)})
	handler := middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
		JWKSClient:        client,
		AuthorizedParties: []string{testParty},
		Clock:             clerktest.NewClockAt(now),
		Leeway:            middleware.DefaultClerkAuthLeeway,
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := serveWithToken(handler, token)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestClerkAuthDiagnosticsClassifyVerificationFailures(t *testing.T) {
	now := time.Date(2026, 7, 12, 20, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		claims          map[string]any
		configure       func(t *testing.T, tokenKey any) []jose.JSONWebKey
		authorizedParty string
		wantCategory    middleware.AuthFailureCategory
	}{
		{
			name: "wrong signature",
			configure: func(t *testing.T, _ any) []jose.JSONWebKey {
				_, otherPublicKey := generateClerkToken(t, "signature-kid", now, nil)
				return []jose.JSONWebKey{newTestJWK("signature-kid", otherPublicKey)}
			},
			wantCategory: middleware.AuthFailureInvalidSignature,
		},
		{
			name: "expired",
			claims: map[string]any{
				"exp": now.Add(-time.Minute).Unix(),
			},
			wantCategory: middleware.AuthFailureTokenExpired,
		},
		{
			name: "not active",
			claims: map[string]any{
				"nbf": now.Add(time.Minute).Unix(),
			},
			wantCategory: middleware.AuthFailureTokenNotActive,
		},
		{
			name: "invalid issuer",
			claims: map[string]any{
				"iss": "https://invalid.example.com",
			},
			wantCategory: middleware.AuthFailureInvalidIssuer,
		},
		{
			name:            "invalid authorized party",
			authorizedParty: "https://other.example.com",
			wantCategory:    middleware.AuthFailureInvalidParty,
		},
	}

	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kid := strings.ReplaceAll(test.name, " ", "-") + "-kid"
			if index == 0 {
				kid = "signature-kid"
			}
			claims := cloneClaims(test.claims)
			if test.authorizedParty != "" {
				claims["azp"] = test.authorizedParty
			}
			token, publicKey := generateClerkToken(t, kid, now, claims)
			keys := []jose.JSONWebKey{newTestJWK(kid, publicKey)}
			if test.configure != nil {
				keys = test.configure(t, publicKey)
			}
			client := newTestJWKSClient(t, keys)
			logs := new(bytes.Buffer)
			logger := slog.New(slog.NewJSONHandler(logs, nil))
			handler := middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
				JWKSClient:        client,
				AuthorizedParties: []string{testParty},
				Clock:             clerktest.NewClockAt(now),
				Logger:            logger,
			})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))

			recorder := serveWithToken(handler, token)
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
			require.Contains(t, logs.String(), `"auth_failure":"`+string(test.wantCategory)+`"`)
			require.NotContains(t, logs.String(), token)
			require.NotContains(t, logs.String(), "user_test")
			require.NotContains(t, logs.String(), "session_test")
			require.NotContains(t, logs.String(), "sk_test_")
		})
	}
}

func TestClerkAuthDiagnosticsClassifyUnknownKeyID(t *testing.T) {
	now := time.Date(2026, 7, 12, 20, 0, 0, 0, time.UTC)
	token, _ := generateClerkToken(t, "missing-kid", now, nil)
	_, otherPublicKey := generateClerkToken(t, "other-kid", now, nil)
	client := newTestJWKSClient(t, []jose.JSONWebKey{newTestJWK("other-kid", otherPublicKey)})
	logs := new(bytes.Buffer)

	handler := middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
		JWKSClient: client,
		Clock:      clerktest.NewClockAt(now),
		Logger:     slog.New(slog.NewJSONHandler(logs, nil)),
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := serveWithToken(handler, token)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Contains(t, logs.String(), `"auth_failure":"unknown_kid"`)
}

func TestNewClerkJWKSClientRejectsPublishableOrMissingKey(t *testing.T) {
	for _, key := range []string{"", "pk_test_public", "wrong"} {
		client, err := middleware.NewClerkJWKSClient(key)
		require.Nil(t, client)
		require.Error(t, err)
		if key != "" {
			require.NotContains(t, err.Error(), key)
		}
	}
}

func TestPreflightClerkJWKS(t *testing.T) {
	_, publicKey := generateClerkToken(t, "preflight-kid", time.Now(), nil)

	t.Run("accepts RS256 signing key", func(t *testing.T) {
		client := newTestJWKSClient(t, []jose.JSONWebKey{newTestJWK("preflight-kid", publicKey)})
		keys, err := middleware.PreflightClerkJWKS(context.Background(), client)
		require.NoError(t, err)
		require.Equal(t, []middleware.ClerkSigningKeyMetadata{{KeyID: "preflight-kid", Algorithm: "RS256"}}, keys)
	})

	t.Run("rejects empty or unsupported set", func(t *testing.T) {
		client := newTestJWKSClient(t, nil)
		keys, err := middleware.PreflightClerkJWKS(context.Background(), client)
		require.Nil(t, keys)
		require.Error(t, err)
	})

	t.Run("reports JWKS endpoint failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		client := jwks.NewClient(&clerk.ClientConfig{BackendConfig: clerk.BackendConfig{
			URL: clerk.String(server.URL),
			Key: clerk.String("sk_test_not-a-secret"),
		}})

		keys, err := middleware.PreflightClerkJWKS(context.Background(), client)
		require.Nil(t, keys)
		require.Error(t, err)
	})
}

func generateClerkToken(t *testing.T, kid string, now time.Time, overrides map[string]any) (string, any) {
	t.Helper()
	claims := map[string]any{
		"iss": testIssuer,
		"sub": "user_test",
		"sid": "session_test",
		"azp": testParty,
		"iat": now.Unix(),
		"nbf": now.Add(-time.Second).Unix(),
		"exp": now.Add(time.Minute).Unix(),
		"v":   2,
	}
	for key, value := range overrides {
		claims[key] = value
	}
	return clerktest.GenerateJWT(t, claims, kid)
}

func cloneClaims(claims map[string]any) map[string]any {
	result := make(map[string]any, len(claims))
	for key, value := range claims {
		result[key] = value
	}
	return result
}

func newTestJWK(kid string, publicKey any) jose.JSONWebKey {
	return jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     kid,
		Algorithm: "RS256",
		Use:       "sig",
	}
}

func newTestJWKSClient(t *testing.T, keys []jose.JSONWebKey) *jwks.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/jwks", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"keys": keys}))
	}))
	t.Cleanup(server.Close)
	return jwks.NewClient(&clerk.ClientConfig{BackendConfig: clerk.BackendConfig{
		URL: clerk.String(server.URL),
		Key: clerk.String("sk_test_not-a-secret"),
	}})
}

func serveWithToken(handler http.Handler, token string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(recorder, request)
	return recorder
}
