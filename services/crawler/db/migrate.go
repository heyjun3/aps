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

func genMigrateInstance(dsn string) (*migrate.Migrate, error) {
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	t, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		panic(err)
	}
	slog.Warn("read dir", "len", len(t))
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		panic(err)
	}
	return migrate.NewWithSourceInstance("iofs",
		d,
		fmt.Sprintf("%s&search_path=%s", dsn, "crawler"),
	)

}

func RunMigrate(dsn string, steps int) {
	m, err := genMigrateInstance(dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Steps(steps); err != nil && err.Error() != "no change" {
		slog.Error("run migrate", "err", err)
	} else {
		slog.Warn("run migrate", "result", "done")
	}
}
