package inventory

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"

	"api-server/database"
	"api-server/spapi"
	"api-server/spapi/price"
)

var SpapiServiceURL string
var db *bun.DB
var inventoryRepository InventoryRepository
var priceRepository PriceRepository
var spapiClient *spapi.SpapiClient

func init() {
	SpapiServiceURL = os.Getenv("SPAPI_SERVICE_URL")
	if SpapiServiceURL == "" {
		panic(errors.New("don't set SPAPI_SERVICE_URL"))
	}
	var err error
	spapiClient, err = spapi.NewSpapiClient(SpapiServiceURL)
	if err != nil {
		panic(err)
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic(errors.New("don't set DB_DSN"))
	}
	db = database.OpenDB(dsn)
	inventoryRepository = InventoryRepository{}
	priceRepository = PriceRepository{}
}

func RefreshInventory(c echo.Context) error {
	var nextToken string
	for {
		res, err := spapiClient.InventorySummaries(nextToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		inventories := Inventories{}
		for _, inventory := range res.Payload.InventorySummaries {
			iv, err := NewInventoryFromInventory(inventory)
			if err != nil {
				slog.Warn(err.Error(), "struct", *inventory, "sku", inventory.SellerSku)
				continue
			}
			inventories = append(inventories, iv)
		}
		if err := inventoryRepository.Save(context.Background(), db, inventories); err != nil {
			slog.Error("failed save inventories", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		slog.Info("logging next token", "nextToken", res.Pagination.NextToken)
		nextToken = res.Pagination.NextToken
		if nextToken == "" {
			slog.Info("break loop")
			break
		}
	}
	return c.JSON(http.StatusOK, "success")
}

func RefreshPricing(c echo.Context) error {
	var inventories Inventories
	var cursor Cursor
	var err error
	for {
		inventories, cursor, err = inventoryRepository.GetNextPage(context.Background(), db, cursor.End, 20)
		if err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		if len(inventories) == 0 || inventories == nil {
			return c.JSON(http.StatusOK, "success")
		}

		res, err := spapiClient.GetPricing(inventories.Skus(), price.Sku)
		if err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		prices := make([]*CurrentPrice, 0, len(inventories))
		for _, payload := range res.Payload {
			offers := payload.Product.Offers
			if len(offers) == 0 {
				continue
			}
			sku := payload.SellerSKU
			amount := offers[0].BuyingPrice.ListingPrice.Amount
			points := offers[0].BuyingPrice.Points.PointsNumber
			if amount == nil || points == nil {
				slog.Warn("amount or points is nil. must be not nil", "sku", sku)
				continue
			}
			price, err := NewCurrentPrice(sku, Ptr(int(*amount)), Ptr(int(*points)))
			if err != nil {
				slog.Error(err.Error(), "sku", sku)
				continue
			}
			prices = append(prices, price)
		}
		if err := priceRepository.Save(context.Background(), db, prices); err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
}
