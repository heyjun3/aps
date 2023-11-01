package spapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

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
	Granularity        Granularity  `json:"granularity"`
	InventorySummaries []*Inventory `json:"inventorySummaries"`
}

type InventorySummariesResponse struct {
	Pagination Pagination `json:"pagination"`
	Payload    Payload    `json:"payload"`
}

func (c SpapiClient) InventorySummaries(nextToken string) (*InventorySummariesResponse, error) {
	query := url.Values{}
	query.Set("nest_token", nextToken)
	c.URL.RawQuery = query.Encode()
	c.URL.Path = "inventory-summaries"

	res, err := http.Get(c.URL.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response InventorySummariesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
