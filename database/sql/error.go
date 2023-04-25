package sql

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsNotFound returns true when the operation that returned an error could not find an item that was expected.
// Applies for Read, ReadValue and ExecOne helpers in this package.
func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows) || errors.Is(err, ErrExecNotOne)
}

func IsForeignKey(err error, constraintName string) bool {
	return violatesConstraint(err, pgerrcode.ForeignKeyViolation, constraintName)
}

func IsUnique(err error, constraintName string) bool {
	return violatesConstraint(err, pgerrcode.UniqueViolation, constraintName)
}

func violatesConstraint(err error, code, name string) bool {
	var e *pgconn.PgError
	if !errors.As(err, &e) {
		return false
	}
	if e.Code != code {
		return false
	}
	return e.ConstraintName == name
}
