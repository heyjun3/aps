package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn), pgdriver.WithTimeout(60*time.Second)))
	return bun.NewDB(sqldb, pgdialect.New())
}

func CreateTable(ctx context.Context, db *bun.DB, model interface{}) error {
	_, err := db.NewCreateTable().
		Model(model).
		IfNotExists().
		Exec(ctx)
	return err
}
