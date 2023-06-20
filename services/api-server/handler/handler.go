package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"

	"api-server/product"
)

var db *bun.DB
func init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db = product.OpenDB(dsn)
}

func GetFilenames(c echo.Context) error {
	ctx := context.Background()
	repo := product.ProductRepository{DB: db}
	filenames, err := repo.GetFilenames(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "error"})
	}
	return c.JSON(http.StatusOK, filenames)
}

func DeleteProducts(c echo.Context) error {
	ctx := context.Background()
	repo := product.ProductRepository{DB: db}
	err := repo.DeleteByFilename(ctx, c.Param("filename"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
