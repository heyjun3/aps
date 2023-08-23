package util

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func DatabaseFactory() (*bun.DB, context.Context) {
	ctx := context.Background()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(TestDSN())))
	conn := bun.NewDB(sqldb, pgdialect.New())
	conn.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return conn, ctx
}

func TestDSN() string {
	DSN := os.Getenv("TEST_DSN")
	if DSN == "" {
		panic("DNS isn't empty")
	}
	return DSN
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

func CreateHttpResponseOnSjis(path string) (*http.Response, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	reader := transform.NewReader(bytes.NewReader(b), japanese.ShiftJIS.NewEncoder())
	res := &http.Response{
		Body:    io.NopCloser(reader),
		Request: &http.Request{},
	}
	return res, nil
}

func OmitError[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
