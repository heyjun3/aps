package handler

import (
	"context"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"

	"api-server/database"
	"api-server/product"
)

type Res struct {
	Message string `json:"status"`
}

type FileListRes struct {
	List []string `json:"list"`
}

type StatusRes struct {
	Keepa map[string]int `json:"keepa"`
	Mws   map[string]int `json:"mws"`
}

type ChartRes struct {
	Charts      []product.ProductWithChart `json:"chart_data"`
	CurrentPage int                        `json:"current_page"`
	MaxPage     int                        `json:"max_page"`
}

var db *bun.DB
var productRepo product.ProductRepository
var keepaRepo product.KeepaRepository

func init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db = database.OpenDB(dsn, true)
	productRepo = product.ProductRepository{DB: db}
	keepaRepo = product.KeepaRepository{DB: db}
}

func GetCharts(c echo.Context) error {
	filename := c.Param("filename")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if filename == "" || err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusBadRequest, Res{"error"})
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	excludeKeywords := slices.DeleteFunc(
		strings.Split(c.QueryParam("excludeKeywords"), ","),
		func(s string) bool {
			return s == ""
		},
	)
	excludeKeywords = append(excludeKeywords, []string{"インク", "トナー"}...)

	ctx := context.Background()
	charts, total, err := productRepo.GetProductWithChartBySearchCondition(
		ctx, product.NewSearchCondition(filename,
			product.SearchConditionWithLimit(limit),
			product.SearchConditionWithPage(page),
			product.SearchConditionWithExcludeKeywords(excludeKeywords),
		),
	)
	if err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}

	maxPage := int(math.Ceil((float64(total) / float64(limit))))
	if page > maxPage {
		slog.Error("page over max page size")
		return c.JSON(http.StatusNotFound, Res{"error"})
	}
	for i := 0; i < len(charts); i++ {
		charts[i].CalculateRankMA(7)
	}

	return c.JSON(http.StatusOK, ChartRes{
		Charts: charts, CurrentPage: page, MaxPage: maxPage})
}

func GetFilenames(c echo.Context) error {
	ctx := context.Background()
	filenames, err := productRepo.GetFilenames(ctx)
	if err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}

	if len(filenames) <= 0 {
		return c.JSON(http.StatusOK, FileListRes{[]string{}})
	}
	return c.JSON(http.StatusOK, FileListRes{filenames})
}

func GetStatusCounts(c echo.Context) error {
	ctx := context.Background()
	keepa, err := keepaRepo.GetCounts(ctx)
	if err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	mws, err := productRepo.GetCounts(ctx)
	if err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, StatusRes{keepa, mws})
}

func DeleteProducts(c echo.Context) error {
	ctx := context.Background()
	err := productRepo.DeleteByFilename(ctx, c.Param("filename"))
	if err != nil {
		slog.Error("error", err)
		return c.JSON(http.StatusInternalServerError, Res{"error"})
	}
	return c.JSON(http.StatusOK, Res{"ok"})
}
