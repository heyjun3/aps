package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

type StatusRes struct {
	Keepa map[string]int `json:"keepa"`
	Mws   map[string]int `json:"mws"`
}

var db *bun.DB
var productRepo product.ProductRepository
var keepaRepo product.KeepaRepository

func init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db = product.OpenDB(dsn)
	productRepo = product.ProductRepository{DB: db}
	keepaRepo = product.KeepaRepository{DB: db}
}

func GetCharts(c echo.Context) error {
	filename := c.Param("filename")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if filename == "" || err != nil {
		return c.JSON(http.StatusBadRequest, Res{"error"})
	}
	
	ctx := context.Background()
	charts, total, err := productRepo.GetProductWithChart(ctx, filename, page, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	
}

func GetFilenames(c echo.Context) error {
	ctx := context.Background()
	filenames, err := productRepo.GetFilenames(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, FileListRes{filenames})
}

func GetStatusCounts(c echo.Context) error {
	ctx := context.Background()
	keepa, err := keepaRepo.GetCounts(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	fmt.Println(keepa)
	mws, err := productRepo.GetCounts(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	fmt.Println(mws)
	return c.JSON(http.StatusOK, StatusRes{keepa, mws})
}

func DeleteProducts(c echo.Context) error {
	ctx := context.Background()
	err := productRepo.DeleteByFilename(ctx, c.Param("filename"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, Res{"ok"})
}
