package inventory

import (
	"context"
	"errors"
	"fmt"
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
	if err := database.CreateTable(context.Background(), db, &Inventory{}); err != nil {
		panic(err)
	}
	inventoryRepository = InventoryRepository{}
}

type Pagination struct {
	NextToken string `json:"nextToken"`
}

type Granularity struct {
	GranularityType string `json:"granularityType"`
	GranularityId   string `json:"granularityId"`
}

type Payload struct {
	Granularity        Granularity  `json:"granularity"`
	InventorySummaries []*Inventory `json:"inventorySummaries"`
}

type InventorySummariesResponse struct {
	Pagination Pagination `json:"pagination"`
	Payload    Payload    `json:"payload"`
}

func RefreshInventory(c echo.Context) error {
	var nextToken string
	for {
		res, err := spapiClient.InventorySummaries(nextToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		inventories := []*Inventory{}
		for _, inventory := range res.Payload.InventorySummaries {
			inventories = append(inventories, &Inventory{Inventory: inventory})
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
	var inventories []*Inventory
	var cursor Cursor
	var err error
	for {
		inventories, cursor, err = inventoryRepository.GetNextPage(context.Background(), db, cursor.End, 20)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if len(inventories) == 0 || inventories == nil {
			return c.JSON(http.StatusOK, "success")
		}
		skus := []string{}
		for _, inventory := range inventories {
			skus = append(skus, inventory.SellerSku)
		}
		res, err := spapiClient.GetPricing(skus, price.Sku)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		for _, payload := range res.Payload {
			fmt.Println("sku", payload.SellerSKU)
			offers := payload.Product.Offers
			if len(offers) == 0 {
				fmt.Println(payload)
			}
			// fmt.Println(offer)
			// fmt.Println("price", payload.Product.Offers[0].BuyingPrice.ListingPrice.Amount)
			// fmt.Println("points", payload.Product.Offers[0].BuyingPrice.Points.PointsNumber)
		}
		return c.JSON(http.StatusOK, "success")
	}
}
