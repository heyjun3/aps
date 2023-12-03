package inventory

import (
	"errors"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"

	"api-server/database"
)

var SpapiServiceURL string
var db *bun.DB
var inventoryService *InventoryService

func init() {
	SpapiServiceURL = os.Getenv("SPAPI_SERVICE_URL")
	if SpapiServiceURL == "" {
		panic(errors.New("don't set SPAPI_SERVICE_URL"))
	}
	var err error
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	db = database.OpenDB(dsn)
	inventoryService, err = NewInventoryService(SpapiServiceURL)
	if err != nil {
		panic(err)
	}
}

func RefreshInventory(c echo.Context) error {
	if err := inventoryService.refreshInventories(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := inventoryService.refreshPricing(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}

func RefreshPricing(c echo.Context) error {
	if err := inventoryService.refreshPricing(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}

func GetInventories(c echo.Context) error {
	inventories, err := inventoryService.getInventories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, inventories)
}

func UpdatePricing(c echo.Context) error {
	dtos := new(updatePricingDTOS)
	if err := c.Bind(dtos); err != nil {
		slog.Error("failed bind body", err)
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	if err := inventoryService.updatePricing(*dtos); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}
