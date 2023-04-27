package common

import (
	"context"
	_ "embed"
	"event-facematch-backend/database/sql"
	"event-facematch-backend/test/integration/testlog"
	"github.com/golang-migrate/migrate/v4"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var database sql.Database

//go:embed scripts/wipe_database.sql
var wipeDatabaseSQL string

//go:embed scripts/populate_database.sql
var populateDatabaseSQL string

// Modified connection string that turns the caching off - it is necessary for testing
var connString = "postgres://postgres:matchtheface123@localhost:5433/event-facematch?sslmode=disable"

type migrateLogger struct{}

func (m migrateLogger) Printf(format string, v ...interface{}) {
	testlog.Logf(format, v...)
}

func (m migrateLogger) Verbose() bool {
	return true
}

func init() {
	ctx := context.Background()

	var err error
	database, err = sql.Open(ctx, connString)
	if err != nil {
		testlog.Logln("Error connecting to database: ", err)
	}
}

// WipePostgres wipes the database.
func WipePostgres() {
	testlog.Logln("WipePostgres: start")
	ctx := context.Background()
	err := sql.WithConnection(ctx, database, func(dctx sql.DataContext) error {
		err := sql.Exec(dctx, wipeDatabaseSQL)

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal("WipePostgres: ", err)
	}

	testlog.Logln("WipePostgres: done")
}

// Process the migration up
func MigrateUp() {
	testlog.Logln("MigrateUp: start")

	m, err := migrate.New("file://database/sql/migrations", connString)
	if err != nil {
		testlog.Logln(err)
		return
	}
	m.Log = migrateLogger{}

	err = m.Up()
	if err != nil {
		testlog.Logln(err)
		return
	}

	testlog.Logln("MigrateUp: done")
}

func PopulatePostgres() {
	testlog.Logln("PopulatePostgres: start")
	ctx := context.Background()
	err := sql.WithConnection(ctx, database, func(dctx sql.DataContext) error {
		err := sql.Exec(dctx, populateDatabaseSQL)

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal("PopulatePostgres: ", err)
	}

	testlog.Logln("PopulatePostgres: done")
}
