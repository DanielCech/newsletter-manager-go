package sql

import (
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func Test_IsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(pgx.ErrNoRows))
}

func Test_IsForeignKey(t *testing.T) {
	constraintName := "user_id_fkey"
	err := &pgconn.PgError{
		Code:           pgerrcode.ForeignKeyViolation,
		ConstraintName: constraintName,
	}
	assert.True(t, IsForeignKey(err, constraintName))
}

func Test_IsUnique(t *testing.T) {
	constraintName := "user_key"
	err := &pgconn.PgError{
		Code:           pgerrcode.UniqueViolation,
		ConstraintName: constraintName,
	}
	assert.True(t, IsUnique(err, constraintName))
}

func Test_violatesConstraint(t *testing.T) {
	constraintName := "user_id_fkey"
	code := pgerrcode.ForeignKeyViolation
	err := &pgconn.PgError{
		Code:           code,
		ConstraintName: constraintName,
	}
	assert.True(t, violatesConstraint(err, code, constraintName))

	err.ConstraintName = "unknown"
	assert.False(t, violatesConstraint(err, code, constraintName))

	assert.False(t, violatesConstraint(err, "unknown_code", constraintName))

	assert.False(t, violatesConstraint(errTest, code, constraintName))
}
