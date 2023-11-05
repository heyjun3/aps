package test

import (
	"context"
	"fmt"
	"os"

	"github.com/uptrace/bun"

	"api-server/database"
)

func CreateTestDBConnection() *bun.DB {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	return database.OpenDB(dsn)
}

func ResetModel(ctx context.Context, db *bun.DB, model interface{}) {
	if err := db.ResetModel(ctx, model); err != nil {
		panic(err)
	}
}
