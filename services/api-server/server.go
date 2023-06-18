package main

import (
	"fmt"
	"context"
	"database/sql"
	// "net/http"

	// "github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"api-server/product"
)

func main() {
	dsn := "postgres://postgres:postgres@postgresql-server:5432/aps?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	repo := product.ProductRepository{DB:db}
	ctx := context.Background()
	p, err := repo.GetFilenames(ctx)
	fmt.Println(err)
	fmt.Println(p)
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.Logger.Fatal(e.Start(":1323"))
}
