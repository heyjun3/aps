package db

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrate(dsn string) {
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs",
		d,
		fmt.Sprintf("%s&search_path=%s", dsn, "crawler"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		slog.Error("run migrate", "err", err)
	} else {
		slog.Warn(err.Error())
	}
}
