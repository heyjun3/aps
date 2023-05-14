package testutil

import (
	"bytes"
	"context"
	"crawler/config"
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func DatabaseFactory() (*bun.DB, context.Context) {
	ctx := context.Background()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(TestDSN())))
	conn := bun.NewDB(sqldb, pgdialect.New())
	return conn, ctx
}

func TestDSN() string {
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	return conf.Dsn()
}

func CreateHttpResponse(path string) (*http.Response, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &http.Response{
		Body:    io.NopCloser(bytes.NewReader(b)),
		Request: &http.Request{},
	}
	return res, nil
}
