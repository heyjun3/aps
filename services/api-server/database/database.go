package database

import (
	"database/sql"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func OpenDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn), pgdriver.WithTimeout(60*time.Second)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(bundebug.WithVerbose(true), bundebug.WithWriter(os.Stdout)),
	)
	return db
}
