package db

import (
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrate(dsn string) {
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("%s&search_path=%s", dsn, "crawler"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		slog.Warn("run migrate", "err", err)
	}
}
