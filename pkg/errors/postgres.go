package errors

import (
	stdErrors "errors"

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
		return NewAppError(CodeUniqueViolation, err, postgresErrorDetails(pgErr, "unique constraint violated")...)
	case PostgresForeignKeyViolation:
		return NewAppError(CodeForeignKey, err, postgresErrorDetails(pgErr, "foreign key constraint violated")...)
	case PostgresCheckViolation:
		return NewAppError(CodeCheckViolation, err, postgresErrorDetails(pgErr, "check constraint violated")...)
	case PostgresNotNullViolation:
		return NewAppError(CodeNotNullViolation, err, postgresErrorDetails(pgErr, "required value is missing")...)
	case PostgresValueTooLong:
		return NewAppError(CodeValueTooLong, err, postgresErrorDetails(pgErr, "value exceeds allowed length")...)
	default:
		return NewAppError(CodeConstraint, err, postgresErrorDetails(pgErr, "database constraint violated")...)
	}
}

func postgresErrorDetails(pgErr *pgconn.PgError, reason string) []ErrorDetail {
	if pgErr == nil {
		return nil
	}

	if pgErr.ColumnName != "" {
		return []ErrorDetail{ConstraintDetail("database.column."+pgErr.ColumnName, reason)}
	}
	if pgErr.ConstraintName != "" {
		return []ErrorDetail{ConstraintDetail("database.constraint."+pgErr.ConstraintName, reason)}
	}
	if pgErr.TableName != "" {
		return []ErrorDetail{ConstraintDetail("database.table."+pgErr.TableName, reason)}
	}

	return nil
}
