package errors

import (
	stdErrors "errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PostgresUniqueViolation     = "23505"
	PostgresForeignKeyViolation = "23503"
	PostgresCheckViolation      = "23514"
	PostgresNotNullViolation    = "23502"
	PostgresValueTooLong        = "22001"
)

func MapPostgresError(err error) *AppError {
	var pgErr *pgconn.PgError
	if !stdErrors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case PostgresUniqueViolation:
		return &AppError{
			Status:  http.StatusConflict,
			Code:    CodeUniqueViolation,
			Message: "resource already exists",
			Err:     err,
		}
	case PostgresForeignKeyViolation:
		return &AppError{
			Status:  http.StatusBadRequest,
			Code:    CodeForeignKey,
			Message: "referenced resource does not exist",
			Err:     err,
		}
	case PostgresCheckViolation:
		return &AppError{
			Status:  http.StatusBadRequest,
			Code:    CodeCheckViolation,
			Message: "request violates a data constraint",
			Err:     err,
		}
	case PostgresNotNullViolation:
		return &AppError{
			Status:  http.StatusBadRequest,
			Code:    CodeNotNullViolation,
			Message: "required value is missing",
			Err:     err,
		}
	case PostgresValueTooLong:
		return &AppError{
			Status:  http.StatusBadRequest,
			Code:    CodeValueTooLong,
			Message: "value is too long",
			Err:     err,
		}
	default:
		return &AppError{
			Status:  http.StatusBadRequest,
			Code:    CodeConstraint,
			Message: "request violates a database constraint",
			Err:     err,
		}
	}
}
