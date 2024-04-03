package main

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	m, err := migrate.New(
		"file://db/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		slog.Warn("run migrate", "err", err)
	}
}
