package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"api-server/handler"
	"api-server/inventory"
	"api-server/shop"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	m, err := migrate.New(
		"file://database/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		slog.Warn("run migrate", "err", err)
	}

}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
	}))
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if c.Request().Method == http.MethodPost {
			fmt.Fprintf(os.Stdout, "Reqest Body %v\n", string(reqBody))
		}
	}))
	validate := validator.New()
	e.Validator = &CustomValidator{validator: shop.AddValidateRules(validate)}

	e.GET("/api/list", handler.GetFilenames)
	e.GET("/api/counts", handler.GetStatusCounts)
	e.GET("/api/chart_list/:filename", handler.GetCharts)
	e.DELETE("/api/deleteFile/:filename", handler.DeleteProducts)

	e.GET("/api/shops", shop.GetShops)
	e.POST("/api/shops", shop.CreateShop)
	e.DELETE("/api/shops", shop.DeleteShop)

	e.GET("/api/inventories", inventory.GetInventories)
	e.POST("/api/inventory/refresh", inventory.RefreshInventory)
	// e.POST("/api/price/refresh", inventory.RefreshPricing)
	e.POST("/api/price/update", inventory.UpdatePricing)
	e.POST("/api/lowest-price/refresh", inventory.RefreshLowestPricing)

	e.Logger.Fatal(e.Start(":5000"))
}
