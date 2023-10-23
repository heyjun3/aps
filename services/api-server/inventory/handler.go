package inventory

import (
	"api-server/database"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"
)

var SpapiServiceURL string
var db *bun.DB
var inventoryRepository InventoryRepository

func init() {
	SpapiServiceURL = os.Getenv("SPAPI_SERVICE_URL")
	if SpapiServiceURL == "" {
		panic(errors.New("don't set SPAPI_SERVICE_URL"))
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
	URL, err := url.Parse(SpapiServiceURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	URL.Path = "inventory-summaries"
	var nextToken string
	for {
		res, err := func() (*InventorySummariesResponse, error) {
			query := url.Values{}
			query.Set("next_token", nextToken)
			URL.RawQuery = query.Encode()
			res, err := http.Get(URL.String())
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			var summariesResponse InventorySummariesResponse
			if err := json.Unmarshal(body, &summariesResponse); err != nil {
				return nil, err
			}
			nextToken = summariesResponse.Pagination.NextToken
			return &summariesResponse, nil
		}()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if err := inventoryRepository.Save(context.Background(), db, res.Payload.InventorySummaries); err != nil {
			slog.Error("failed save inventories", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		if nextToken == "" {
			slog.Info("break loop")
			break
		}
	}
	return c.JSON(http.StatusOK, "success")
}
