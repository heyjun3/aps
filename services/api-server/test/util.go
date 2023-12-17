package test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/uptrace/bun"

	"api-server/database"
)

var dsn string

func init() {
	dsn = os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	m, err := migrate.New(
		"file://../database/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		fmt.Println(err)
	}
}

func CreateTestDBConnection() *bun.DB {
	return database.OpenDB(dsn, true)
}

func ResetModel(ctx context.Context, db *bun.DB, model interface{}) {
	if err := db.ResetModel(ctx, model); err != nil {
		panic(err)
	}
}

func OmitErr[T any](v *T, err error) *T {
	if err != nil {
		panic(err)
	}
	return v
}
