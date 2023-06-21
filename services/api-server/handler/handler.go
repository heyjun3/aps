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

type Res struct {
	Message string `json:"message"`
}

type FileListRes struct {
	List []string `json:"list"`
}

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
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, FileListRes{filenames})
}

func DeleteProducts(c echo.Context) error {
	ctx := context.Background()
	repo := product.ProductRepository{DB: db}
	err := repo.DeleteByFilename(ctx, c.Param("filename"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, Res{"ok"})
}
