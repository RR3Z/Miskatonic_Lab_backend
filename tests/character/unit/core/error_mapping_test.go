package tests

import (
	"errors"
	"net/http"
	"testing"

	characterHandlerErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character/errors"
	characterServiceErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCharacterHandlerMapsGenericDomainErrors(t *testing.T) {
	tests := []struct {
		name       string
		appErr     func() errorMappingResult
		wantStatus int
		wantCode   string
	}{
		{
			name: "invalid path id",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.InvalidPathIDError("invalid note id", errors.New("invalid uuid"))
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.invalid_id",
		},
		{
			name: "invalid input",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.InvalidInputError("invalid request body", errors.New("bad json"))
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.invalid_input",
		},
		{
			name: "not found",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.MapNotFoundOrServiceError(pgx.ErrNoRows, "character not found", "failed")
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "character.not_found",
		},
		{
			name: "age negative",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.MapServiceError(characterServiceErrors.ErrAgeNegative, "failed")
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.age_negative",
		},
		{
			name: "sex invalid",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.MapServiceError(characterServiceErrors.ErrSexInvalid, "failed")
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "character.sex_invalid",
		},
		{
			name: "character limit reached",
			appErr: func() errorMappingResult {
				err := characterHandlerErrors.MapServiceError(characterServiceErrors.ErrCharacterLimitReached, "failed")
				return errorMappingResult{status: err.StatusCode(), code: err.Response().Code}
			},
			wantStatus: http.StatusConflict,
			wantCode:   "character.limit_reached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appErr()

			require.Equal(t, tt.wantStatus, got.status)
			require.Equal(t, tt.wantCode, got.code)
		})
	}
}

func TestCharacterLimitErrorIncludesConflictDetail(t *testing.T) {
	appErr := characterHandlerErrors.MapServiceError(characterServiceErrors.ErrCharacterLimitReached, "failed")
	response := appErr.Response()

	require.Equal(t, http.StatusConflict, appErr.StatusCode())
	require.Equal(t, "character.limit_reached", response.Code)
	require.Len(t, response.Details, 1)
	require.Equal(t, "conflict", response.Details[0].Type)
	require.Equal(t, "characters", response.Details[0].Target)
	require.Equal(t, "limit_reached", response.Details[0].Reason)
}

type errorMappingResult struct {
	status int
	code   string
}
