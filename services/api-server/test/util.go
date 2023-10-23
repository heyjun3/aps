package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func CreateTestDBConnection() *bun.DB {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(bundebug.WithVerbose(true), bundebug.WithWriter(os.Stdout)),
	)

	return db
}

func ResetModel(ctx context.Context, db *bun.DB, model interface{}) {
	if err := db.ResetModel(ctx, model); err != nil {
		panic(err)
	}
}
