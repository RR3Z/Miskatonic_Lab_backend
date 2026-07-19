package errors

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed error_catalog.json
var errorCatalogJSON []byte

type ErrorDefinition struct {
	HTTPStatus         int      `json:"httpStatus"`
	Message            string   `json:"message"`
	Category           string   `json:"category"`
	AllowedDetailTypes []string `json:"allowedDetailTypes"`
}

type errorCatalog struct {
	Version int                        `json:"version"`
	Errors  map[string]ErrorDefinition `json:"errors"`
}

var catalog errorCatalog

func init() {
	if err := json.Unmarshal(errorCatalogJSON, &catalog); err != nil {
		panic(fmt.Sprintf("decode error catalog: %v", err))
	}
	if err := validateCatalog(catalog); err != nil {
		panic(fmt.Sprintf("validate error catalog: %v", err))
	}
}

func ErrorDefinitionFor(code string) (ErrorDefinition, bool) {
	definition, ok := catalog.Errors[code]
	return definition, ok
}

func ErrorCodes() []string {
	codes := make([]string, 0, len(catalog.Errors))
	for code := range catalog.Errors {
		codes = append(codes, code)
	}
	return codes
}

func NewAppError(code string, err error, details ...ErrorDetail) *AppError {
	definition, ok := ErrorDefinitionFor(code)
	if !ok {
		definition, _ = ErrorDefinitionFor(CodeInternalError)
		return &AppError{
			Status:  definition.HTTPStatus,
			Code:    CodeInternalError,
			Message: definition.Message,
			Err:     fmt.Errorf("unknown error code %q: %w", code, err),
		}
	}

	return &AppError{
		Status:  definition.HTTPStatus,
		Code:    code,
		Message: definition.Message,
		Details: details,
		Err:     err,
	}
}

func validateCatalog(value errorCatalog) error {
	if value.Version != 1 || len(value.Errors) == 0 {
		return fmt.Errorf("version must be 1 and errors must be non-empty")
	}
	allowedCategories := map[string]bool{
		"validation": true, "parse": true, "permission": true, "resource_state": true,
		"conflict": true, "constraint": true, "internal": true,
	}
	allowedDetails := map[string]bool{
		"validation": true, "parse": true, "resource_state": true, "permission": true,
		"constraint": true, "conflict": true,
	}
	for code, definition := range value.Errors {
		if !strings.Contains(code, ".") || strings.HasPrefix(code, ".") || strings.HasSuffix(code, ".") {
			return fmt.Errorf("invalid code %q", code)
		}
		if definition.HTTPStatus < 400 || definition.HTTPStatus > 599 {
			return fmt.Errorf("%s: invalid HTTP status", code)
		}
		if definition.Message == "" || !allowedCategories[definition.Category] {
			return fmt.Errorf("%s: message or category invalid", code)
		}
		for _, detailType := range definition.AllowedDetailTypes {
			if !allowedDetails[detailType] {
				return fmt.Errorf("%s: invalid detail type %q", code, detailType)
			}
		}
	}
	if _, ok := value.Errors[CodeInternalError]; !ok {
		return fmt.Errorf("%s is required", CodeInternalError)
	}
	return nil
}
