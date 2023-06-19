package main

import (
	"fmt"
	"context"
	// "net/http"

	// "github.com/labstack/echo/v4"

	"api-server/product"
)



func main() {
	dsn := "postgres://postgres:postgres@postgresql-server:5432/aps?sslmode=disable"
	db := product.OpenDB(dsn)
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
