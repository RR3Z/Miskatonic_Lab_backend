package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorCatalogDefinitionsAreUsable(t *testing.T) {
	for _, code := range ErrorCodes() {
		definition, ok := ErrorDefinitionFor(code)
		if !ok {
			t.Fatalf("missing definition for %s", code)
		}
		if definition.HTTPStatus < http.StatusBadRequest || definition.HTTPStatus > 599 {
			t.Fatalf("invalid HTTP status for %s: %d", code, definition.HTTPStatus)
		}

		appErr := NewAppError(code, errors.New("cause"))
		if appErr.Code != code || appErr.Status != definition.HTTPStatus || appErr.Message != definition.Message {
			t.Fatalf("factory did not use catalog definition for %s", code)
		}
	}
}

func TestNewAppErrorFallsBackForUnknownCode(t *testing.T) {
	appErr := NewAppError("missing.code", errors.New("cause"))
	definition, _ := ErrorDefinitionFor(CodeInternalError)

	if appErr.Code != CodeInternalError || appErr.Status != definition.HTTPStatus || appErr.Message != definition.Message {
		t.Fatalf("unexpected fallback: %#v", appErr)
	}
}

func TestNormalizeAppErrorUsesCatalogButKeepsDetails(t *testing.T) {
	detail := ValidationDetail("body.name", "required")
	appErr := NormalizeAppError(&AppError{
		Status:  http.StatusTeapot,
		Code:    "character.name_required",
		Message: "stale message",
		Details: []ErrorDetail{detail},
	})

	if appErr.Status != http.StatusBadRequest || appErr.Message != "name is required" {
		t.Fatalf("catalog was not applied: %#v", appErr)
	}
	if len(appErr.Details) != 1 || appErr.Details[0] != detail {
		t.Fatalf("details were changed: %#v", appErr.Details)
	}
}
