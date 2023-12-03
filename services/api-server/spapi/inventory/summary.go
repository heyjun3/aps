package inventory

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type SummariesResponse struct {
	Pagination Pagination `json:"pagination"`
	Payload    Payload    `json:"payload"`
}

type Pagination struct {
	NextToken string `json:"nextToken"`
}

type Payload struct {
	Granularity        Granularity  `json:"granularity"`
	InventorySummaries []*Inventory `json:"inventorySummaries"`
}

type Granularity struct {
	GranularityType string `json:"granularityType"`
	GranularityId   string `json:"granularityId"`
}

type Inventory struct {
	Asin            *string `json:"asin" bun:"asin"`
	FnSku           *string `json:"fnSku" bun:"fnsku"`
	SellerSku       *string `json:"sellerSku" bun:"seller_sku,pk"`
	Condition       *string `json:"condition" bun:"condition"`
	LastUpdatedTime *string `json:"lastUpdatedTime" bun:"-"`
	ProductName     *string `json:"productName" bun:"product_name"`
	TotalQuantity   *int    `json:"totalQuantity" bun:"quantity"`
}

func Summaries(URL *url.URL, nextToken string) (*SummariesResponse, error) {
	query := url.Values{}
	query.Set("next_token", nextToken)
	URL.RawQuery = query.Encode()
	URL.Path = "inventory-summaries"

	res, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response SummariesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
