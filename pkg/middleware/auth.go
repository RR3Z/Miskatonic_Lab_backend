package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	josejwt "github.com/go-jose/go-jose/v3/jwt"
)

type ClerkAuthConfig struct {
	JWKSClient        *jwks.Client
	AuthorizedParties []string
	Logger            *slog.Logger
	Clock             clerk.Clock
	Leeway            time.Duration
}

type ClerkSigningKeyMetadata struct {
	KeyID     string
	Algorithm string
}

type AuthFailureCategory string

const (
	DefaultClerkAuthLeeway                            = 30 * time.Second
	AuthFailureMissingOrMalformed AuthFailureCategory = "missing_or_malformed_token"
	AuthFailureJWKFetchFailed     AuthFailureCategory = "jwk_fetch_failed"
	AuthFailureUnknownKeyID       AuthFailureCategory = "unknown_kid"
	AuthFailureInvalidSignature   AuthFailureCategory = "invalid_signature"
	AuthFailureTokenExpired       AuthFailureCategory = "token_expired"
	AuthFailureTokenNotActive     AuthFailureCategory = "token_not_active"
	AuthFailureInvalidIssuer      AuthFailureCategory = "invalid_issuer"
	AuthFailureInvalidParty       AuthFailureCategory = "invalid_authorized_party"
)

type tokenMetadata struct {
	Algorithm       string `json:"alg"`
	KeyID           string `json:"kid"`
	Issuer          string `json:"iss"`
	AuthorizedParty string `json:"azp"`
	ExpiresAt       int64  `json:"exp"`
	NotBefore       int64  `json:"nbf"`
}

func NewClerkJWKSClient(secretKey string) (*jwks.Client, error) {
	if !strings.HasPrefix(secretKey, "sk_test_") && !strings.HasPrefix(secretKey, "sk_live_") {
		return nil, fmt.Errorf("CLERK_SECRET_KEY must start with sk_test_ or sk_live_")
	}

	config := &clerk.ClientConfig{}
	config.Key = clerk.String(secretKey)
	return jwks.NewClient(config), nil
}

func PreflightClerkJWKS(ctx context.Context, client *jwks.Client) ([]ClerkSigningKeyMetadata, error) {
	if client == nil {
		return nil, fmt.Errorf("clerk JWKS client is nil")
	}

	keySet, err := client.Get(ctx, &jwks.GetParams{})
	if err != nil {
		return nil, fmt.Errorf("fetch Clerk JWKS: %w", err)
	}

	keys := make([]ClerkSigningKeyMetadata, 0, len(keySet.Keys))
	for _, key := range keySet.Keys {
		if key == nil || key.KeyID == "" || key.Algorithm != "RS256" {
			continue
		}
		keys = append(keys, ClerkSigningKeyMetadata{
			KeyID:     key.KeyID,
			Algorithm: key.Algorithm,
		})
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("Clerk JWKS contains no RS256 signing keys")
	}

	return keys, nil
}

func NewClerkAuthMiddleware(config ClerkAuthConfig) func(http.Handler) http.Handler {
	logger := config.Logger
	if logger == nil {
		logger = slog.Default()
	}

	options := make([]clerkhttp.AuthorizationOption, 0, 5)
	if config.JWKSClient != nil {
		options = append(options, clerkhttp.JWKSClient(config.JWKSClient))
	}
	if len(config.AuthorizedParties) > 0 {
		options = append(options, clerkhttp.AuthorizedPartyMatches(config.AuthorizedParties...))
	}
	if config.Clock != nil {
		options = append(options, clerkhttp.Clock(config.Clock))
	}
	if config.Leeway > 0 {
		options = append(options, clerkhttp.Leeway(config.Leeway))
	}
	options = append(options, clerkhttp.AuthorizationFailureHandler(
		newClerkAuthFailureHandler(config, logger),
	))

	clerkMiddleware := clerkhttp.RequireHeaderAuthorization(options...)
	return func(next http.Handler) http.Handler {
		protected := clerkMiddleware(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			metadata := decodeTokenMetadata(token)
			if token == "" || metadata.KeyID == "" {
				logClerkAuthFailure(r, logger, AuthFailureMissingOrMalformed, metadata, config.Clock)
			}
			protected.ServeHTTP(w, r)
		})
	}
}

func newClerkAuthFailureHandler(config ClerkAuthConfig, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		metadata := decodeTokenMetadata(token)
		category := diagnoseClerkAuthFailure(r.Context(), token, metadata, config)
		logClerkAuthFailure(r, logger, category, metadata, config.Clock)
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func logClerkAuthFailure(r *http.Request, logger *slog.Logger, category AuthFailureCategory, metadata tokenMetadata, clock clerk.Clock) {
	now := time.Now().UTC()
	if clock != nil {
		now = clock.Now().UTC()
	}

	logger.WarnContext(r.Context(), "clerk authorization failed",
		"auth_failure", category,
		"method", r.Method,
		"path", r.URL.Path,
		"alg", metadata.Algorithm,
		"kid", metadata.KeyID,
		"iss", metadata.Issuer,
		"azp", metadata.AuthorizedParty,
		"exp", metadata.ExpiresAt,
		"nbf", metadata.NotBefore,
		"server_time", now.Unix(),
	)
}

func diagnoseClerkAuthFailure(ctx context.Context, token string, metadata tokenMetadata, config ClerkAuthConfig) AuthFailureCategory {
	if token == "" || metadata.KeyID == "" {
		return AuthFailureMissingOrMalformed
	}

	jwk, err := clerkjwt.GetJSONWebKey(ctx, &clerkjwt.GetJSONWebKeyParams{
		KeyID:      metadata.KeyID,
		JWKSClient: config.JWKSClient,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "missing json web key") {
			return AuthFailureUnknownKeyID
		}
		return AuthFailureJWKFetchFailed
	}

	params := &clerkjwt.VerifyParams{
		Token:  token,
		JWK:    jwk,
		Clock:  config.Clock,
		Leeway: config.Leeway,
	}
	if len(config.AuthorizedParties) > 0 {
		allowed := make(map[string]struct{}, len(config.AuthorizedParties))
		for _, party := range config.AuthorizedParties {
			allowed[party] = struct{}{}
		}
		params.AuthorizedPartyHandler = func(party string) bool {
			if party == "" {
				return true
			}
			_, ok := allowed[party]
			return ok
		}
	}

	_, err = clerkjwt.Verify(ctx, params)
	if err == nil {
		return AuthFailureJWKFetchFailed
	}
	return classifyClerkVerificationError(err)
}

func classifyClerkVerificationError(err error) AuthFailureCategory {
	switch {
	case errors.Is(err, josejwt.ErrExpired):
		return AuthFailureTokenExpired
	case errors.Is(err, josejwt.ErrNotValidYet):
		return AuthFailureTokenNotActive
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "invalid issuer"):
		return AuthFailureInvalidIssuer
	case strings.Contains(message, "invalid authorized party"):
		return AuthFailureInvalidParty
	case strings.Contains(message, "signature"), strings.Contains(message, "verification error"):
		return AuthFailureInvalidSignature
	default:
		return AuthFailureInvalidSignature
	}
}

func extractBearerToken(r *http.Request) string {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(authorization, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
}

func decodeTokenMetadata(token string) tokenMetadata {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenMetadata{}
	}

	var header tokenMetadata
	if !decodeJWTPart(parts[0], &header) {
		return tokenMetadata{}
	}
	var claims tokenMetadata
	if !decodeJWTPart(parts[1], &claims) {
		return tokenMetadata{}
	}

	header.Issuer = claims.Issuer
	header.AuthorizedParty = claims.AuthorizedParty
	header.ExpiresAt = claims.ExpiresAt
	header.NotBefore = claims.NotBefore
	return header
}

func decodeJWTPart(part string, target any) bool {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return false
	}
	return json.Unmarshal(data, target) == nil
}
