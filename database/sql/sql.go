package sql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultSSLMode = "disable"
)

var (
	ErrNoQuerier = errors.New("no querier found")
	ErrNoConn    = errors.New("no connection found")
	ErrNoTx      = errors.New("no transaction found")
)

// Config contains fields required for creating a pgx connection pool.
type Config struct {
	Secret Secret `json:"secret" yaml:"secret" env:",dive"`
}

// Secret contains Path to file or ARN of secrets manager object where are stored config values.
type Secret struct {
	Path *string `json:"path" yaml:"path" env:"DATABASE_SECRET_PATH"`
	ARN  *string `json:"arn" yaml:"arn" env:"DATABASE_SECRET_ARN"`
}

// SecretARN returns ARN of secrets manager object.
func (s Secret) SecretARN() *string {
	return s.ARN
}

// LocalPath returns local path to secret.
func (s Secret) LocalPath() *string {
	return s.Path
}

// DSNValues contains fields which are needed to build data source name.
type DSNValues struct {
	Host     string      `json:"host" validate:"required"`
	Port     json.Number `json:"port" validate:"required"`
	UserName string      `json:"username" validate:"required"`
	Password string      `json:"password" validate:"required"`
	DBName   string      `json:"dbname" validate:"required"`
	SSLMode  string      `json:"sslmode"`
}

// ConnString returns connection string for a database.
func (d DSNValues) ConnString() (string, error) {
	port, err := d.Port.Int64()
	if err != nil {
		return "", fmt.Errorf("parsing port: %w", err)
	}
	if d.SSLMode == "" {
		d.SSLMode = defaultSSLMode
	}

	u := url.URL{}
	u.Scheme = "postgres"
	u.User = url.UserPassword(d.UserName, d.Password)
	u.Host = fmt.Sprintf("%s:%d", d.Host, port)
	u.Path = "/" + d.DBName
	query := u.Query()
	query.Set("sslmode", d.SSLMode)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

type PGXPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// Database contains pgx connection pool.
type Database struct {
	pool PGXPool
}

// Open returns initialized database with a connection pool.
func Open(ctx context.Context, dsn string) (Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return Database{}, err
	}
	return Database{
		pool: pool,
	}, nil
}

// Close closes database with connection pool.
func (d Database) Close() {
	d.pool.Close()
}

// AcquireConnCtx returns context containing pgx connection.
func (d Database) AcquireConnCtx(ctx context.Context) (DataContext, error) {
	c, err := d.pool.Acquire(ctx)
	if err != nil {
		return dbContext{}, err
	}
	return dbContext{Context: WithQuerier(ctx, c)}, nil
}

// ReleaseConnCtx releases a connection that is found in the context.
func (Database) ReleaseConnCtx(dctx DataContext) error {
	querier, ok := QuerierFromCtx(dctx.Ctx())
	if !ok {
		return ErrNoConn
	}
	c, ok := querier.(*pgxpool.Conn)
	if !ok {
		return ErrNoConn
	}
	c.Release()
	return nil
}

// Begin initializes a pgx transaction that is returned in the context.
func (d Database) Begin(ctx context.Context) (DataContext, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return dbContext{Context: WithQuerier(ctx, tx)}, nil
}

// Commit commits a pgx transaction found in the context.
func (Database) Commit(dctx DataContext) error {
	querier, ok := QuerierFromCtx(dctx.Ctx())
	if !ok {
		return ErrNoTx
	}
	tx, ok := querier.(pgx.Tx)
	if !ok {
		return ErrNoTx
	}
	if err := tx.Commit(dctx.Ctx()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return err
	}
	return nil
}

// Rollback rollbacks a pgx transaction found in the context.
func (Database) Rollback(dctx DataContext) error {
	querier, ok := QuerierFromCtx(dctx.Ctx())
	if !ok {
		return ErrNoTx
	}
	tx, ok := querier.(pgx.Tx)
	if !ok {
		return ErrNoTx
	}
	if err := tx.Rollback(dctx.Ctx()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return err
	}
	return nil
}

// Querier contains functions for querying a database.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func querier(dctx DataContext) Querier {
	q, ok := dctx.Ctx().Value(contextKey.querier).(Querier)
	if !ok {
		panic(ErrNoQuerier)
	}
	return q
}
