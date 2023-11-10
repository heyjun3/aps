package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
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

	e.POST("/api/refresh-inventory", inventory.RefreshInventory)
	// e.POST("/api/refresh-pricing", inventory.RefreshPricing)

	e.Logger.Fatal(e.Start(":5000"))
}
