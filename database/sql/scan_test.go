package sql

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Read(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	query := `SELECT name, email FROM "user" WHERE id = @id`
	type user struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	args := pgx.NamedArgs{"id": 1}
	expected := &user{
		Name:  "Jozko Dlouhy",
		Email: "jozko.dlouhy@gmail.com",
	}
	row := pgxmock.NewRows([]string{"name", "email"}).AddRow(expected.Name, expected.Email)
	mockPool.ExpectQuery(query).WithArgs(args).WillReturnRows(row)

	u, err := Read[user](dctx, query, args)
	assert.NoError(t, err)
	assert.Equal(t, expected, u)

	mockPool.ExpectQuery(query).WithArgs(args).WillReturnError(errTest)
	u, err = Read[user](dctx, query, args)
	assert.Error(t, err)
	assert.Empty(t, u)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_ReadValue(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	query := `SELECT name, email FROM "user" WHERE id = @id`
	type user struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	args := pgx.NamedArgs{"id": 1}
	expected := user{
		Name:  "Jozko Dlouhy",
		Email: "jozko.dlouhy@gmail.com",
	}
	row := pgxmock.NewRows([]string{"name", "email"}).AddRow(expected.Name, expected.Email)
	mockPool.ExpectQuery(query).WithArgs(args).WillReturnRows(row)

	u, err := ReadValue[user](dctx, query, args)
	assert.NoError(t, err)
	assert.Equal(t, expected, u)

	mockPool.ExpectQuery(query).WithArgs(args).WillReturnError(errTest)
	u, err = ReadValue[user](dctx, query, args)
	assert.Error(t, err)
	assert.Empty(t, u)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_List(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	query := `SELECT name, email FROM "user" WHERE created_at > @created_at`
	type user struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	args := pgx.NamedArgs{"created_at": time.Now()}
	expected := []user{
		{
			Name:  "Jozko Dlouhy",
			Email: "jozko.dlouhy@gmail.com",
		},
		{
			Name:  "Karel Klatil",
			Email: "karel.klatil@gmail.com",
		},
	}
	rows := pgxmock.NewRows([]string{"name", "email"}).AddRows(
		[]any{expected[0].Name, expected[0].Email},
		[]any{expected[1].Name, expected[1].Email},
	)
	mockPool.ExpectQuery(query).WithArgs(args).WillReturnRows(rows)

	u, err := List[user](dctx, query, args)
	assert.NoError(t, err)
	assert.Equal(t, expected, u)

	mockPool.ExpectQuery(query).WithArgs(args).WillReturnError(errTest)
	u, err = List[user](dctx, query, args)
	assert.Error(t, err)
	assert.Empty(t, u)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_ExecOne(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	query := `DELETE FROM "user" WHERE id > @id`
	args := pgx.NamedArgs{"id": 1}
	mockPool.ExpectExec(query).WithArgs(args).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = ExecOne(dctx, query, args)
	assert.NoError(t, err)

	mockPool.ExpectExec(query).WithArgs(args).WillReturnResult(pgxmock.NewResult("DELETE", 6))
	err = ExecOne(dctx, query, args)
	assert.ErrorIs(t, err, ErrExecNotOne)

	mockPool.ExpectExec(query).WithArgs(args).WillReturnError(errTest)
	err = ExecOne(dctx, query, args)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrExecNotOne)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_Exec(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	query := `DELETE FROM "user" WHERE id > @id`
	args := pgx.NamedArgs{"id": 1}
	mockPool.ExpectExec(query).WithArgs(args).WillReturnResult(pgxmock.NewResult("DELETE", 3))

	err = Exec(dctx, query, args)
	assert.NoError(t, err)

	mockPool.ExpectExec(query).WithArgs(args).WillReturnError(errTest)
	err = Exec(dctx, query, args)
	assert.Error(t, err)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}
