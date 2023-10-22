package inventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

var SpapiServiceURL string

func init() {
	SpapiServiceURL = os.Getenv("SPAPI_SERVICE_URL")
	if SpapiServiceURL == "" {
		panic(errors.New("don't set SPAPI_SERVICE_URL"))
	}
}

type Pagination struct {
	NextToken string `json:"nextToken"`
}

type Granularity struct {
	GranularityType string `json:"granularityType"`
	GranularityId   string `json:"granularityId"`
}

type Inventory struct {
	Asin            string `json:"asin"`
	FnSku           string `json:"fnSku"`
	SellerSku       string `json:"sellerSku"`
	Condition       string `json:"condition"`
	LastUpdatedTime string `json:"lastUpdatedTime"`
	ProductName     string `json:"productName"`
	TotalQuantity   int    `json:"totalQuantity"`
}

type Payload struct {
	Granularity        Granularity `json:"granularity"`
	InventorySummaries []Inventory `json:"inventorySummaries"`
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
		fmt.Println(res.Payload.InventorySummaries)
		if nextToken == "" {
			slog.Info("break loop")
			break
		}
	}
	return c.JSON(http.StatusOK, "success")
}
