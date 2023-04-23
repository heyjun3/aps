package testutil

import (
	"context"
	"crawler/config"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func DatabaseFactory() (*bun.DB, context.Context) {
	ctx := context.Background()
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.Dsn())))
	conn := bun.NewDB(sqldb, pgdialect.New())
	return conn, ctx
}
