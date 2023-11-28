package inventory

import (
	"context"
	"errors"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
var priceRepository PriceRepository[*CurrentPrice]
var lowestPriceRepository PriceRepository[*LowestPrice]
var desiredPriceRepository PriceRepository[*DesiredPrice]
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
	priceRepository = PriceRepository[*CurrentPrice]{}
	lowestPriceRepository = PriceRepository[*LowestPrice]{}
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

// func RefreshPricing(c echo.Context) error {
// 	var inventories Inventories
// 	var cursor Cursor
// 	var err error
// 	for {
// 		inventories, cursor, err = inventoryRepository.GetNextPage(context.Background(), db, cursor.End, 20)
// 		if err != nil {
// 			slog.Error("error", "detail", err)
// 			return c.JSON(http.StatusInternalServerError, err)
// 		}
// 		if len(inventories) == 0 || inventories == nil {
// 			return c.JSON(http.StatusOK, "success")
// 		}

// 		res, err := spapiClient.GetPricing(inventories.Skus(), price.Sku)
// 		if err != nil {
// 			slog.Error("error", "detail", err)
// 			return c.JSON(http.StatusInternalServerError, err)
// 		}

// 		prices := make(CurrentPrices, 0, len(inventories))
// 		for _, payload := range res.Payload {
// 			offers := payload.Product.Offers
// 			if len(offers) == 0 {
// 				continue
// 			}
// 			sku := payload.SellerSKU
// 			amount := offers[0].BuyingPrice.ListingPrice.Amount
// 			points := offers[0].BuyingPrice.Points.PointsNumber
// 			if amount == nil || points == nil {
// 				slog.Warn("amount or points is nil. must be not nil", "sku", sku)
// 				continue
// 			}
// 			price, err := NewCurrentPrice(sku, Ptr(int(*amount)), Ptr(int(*points)))
// 			if err != nil {
// 				slog.Error(err.Error(), "sku", sku)
// 				continue
// 			}
// 			prices = append(prices, price)
// 		}
// 		if err := priceRepository.Save(context.Background(), db, prices); err != nil {
// 			slog.Error("error", "detail", err)
// 			return c.JSON(http.StatusInternalServerError, err)
// 		}
// 	}
// }

func RefreshLowestPricing(c echo.Context) error {
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

		res, err := spapiClient.GetPricing(inventories.Skus())
		if err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		currentPrices := make(CurrentPrices, 0, len(inventories))
		lowestPrices := make(LowestPrices, 0, len(inventories))
		for _, response := range res.Responses {
			offers := response.Body.Payload.Offers
			if len(offers) == 0 {
				continue
			}
			sku := response.Body.Payload.SKU
			currents := offers.FilterCondition(price.Condition{MyOffer: Ptr(true)})
			if len(currents) > 0 {
				amount := currents[0].Price.Amount
				point := currents[0].Points.PointsNumber
				price, err := NewCurrentPrice(*sku, *amount, point)
				if err == nil {
					currentPrices = append(currentPrices, price)
				}
			}

			lowest := getLowestPricingFromOffers(offers)
			if lowest != nil {
				amount := lowest.Price.Amount
				point := lowest.Points.PointsNumber
				price, err := NewLowestPrice(*sku, *amount, point)
				if err != nil {
					slog.Error(err.Error(), "sku", sku)
					continue
				}
				lowestPrices = append(lowestPrices, price)
			}
		}
		if err := lowestPriceRepository.Save(context.Background(), db, lowestPrices); err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		if err := priceRepository.Save(context.Background(), db, currentPrices); err != nil {
			slog.Error("error", "detail", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
}

func getLowestPricingFromOffers(offers price.Offers) *price.Offer {
	buyBoxWinnerAndFullfilledByAmazon := offers.FilterCondition(price.Condition{
		IsFullfilledByAmazon: Ptr(true),
		IsBuyBoxWinner:       Ptr(true),
	})
	if len(buyBoxWinnerAndFullfilledByAmazon) > 0 {
		return &buyBoxWinnerAndFullfilledByAmazon[0]
	}
	buyBox := offers.FilterCondition(price.Condition{
		IsBuyBoxWinner: Ptr(true),
	})
	if len(buyBox) > 0 {
		return &buyBox[0]
	}
	return offers.Lowest()
}

func GetInventories(c echo.Context) error {
	ctx := context.Background()
	condition := Condition{Quantity: Ptr(0), IsNotOnlyLowestPrice: true}
	inventories, err := inventoryRepository.GetByCondition(ctx, db, condition)
	if err != nil {
		slog.Error("get inventories error", "detail", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, inventories)
}

type UpdatePricingDTO struct {
	Sku          string `json:"sku"`
	Price        int    `json:"price"`
	PercnetPoint int    `json:"percentPoint"`
}

func UpdatePricing(c echo.Context) error {
	dtos := new([]UpdatePricingDTO)
	if err := c.Bind(dtos); err != nil {
		slog.Error("failed bind body", err)
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	skus := make([]string, 0, len(*dtos))
	for _, dto := range *dtos {
		d := dto
		skus = append(skus, d.Sku)
	}
	condition := Condition{Skus: skus}
	inventories, err := inventoryRepository.GetByCondition(context.Background(), db, condition)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	m := inventories.Map()

	prices := make(DesiredPrices, 0, len(*dtos))
	for _, dto := range *dtos {
		d := dto
		inventory := m[d.Sku]
		if inventory == nil {
			continue
		}
		lowestPrice := inventory.LowestPrice
		if lowestPrice == nil {
			continue
		}

		price, err := NewDesiredPrice(&d.Sku, &d.Price, &d.PercnetPoint, *lowestPrice)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		prices = append(prices, price)
	}
	if err := desiredPriceRepository.Save(context.Background(), db, prices); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, price := range prices {
		p := price
		if err := spapiClient.UpdatePricing(p); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	if err := spapiClient.UpdatePoints(prices); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, dtos)
}
