package lowest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/slog"
)

type GetLowestPricingResponse struct {
	Responses []Response `json:"responses"`
}
type Response struct {
	Body Body `json:"body"`
}
type Body struct {
	Payload Payload `json:"payload"`
}
type Payload struct {
	Offers []Offer `json:"Offers"`
	SKU    *string `json:"SKU"`
}
type Offer struct {
	Points Point `json:"Points"`
	Price  Price `json:"ListingPrice"`
}
type Point struct {
	PointsNumber *int64 `json:"PointsNumber"`
}
type Price struct {
	Amount *float64 `json:"Amount"`
}

func GetLowestPricing(URL *url.URL, skus []string) (*GetLowestPricingResponse, error) {
	if len(skus) == 0 {
		return nil, errors.New("skus must contain at least one letter")
	}

	query := url.Values{}
	query.Set("skus", strings.Join(skus, ","))
	URL.Path = "get-listing-offers-batch"
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
	var getLowestPricingResponse GetLowestPricingResponse
	if err := json.Unmarshal(body, &getLowestPricingResponse); err != nil {
		slog.Error("err", err)
		return nil, err
	}
	return &getLowestPricingResponse, nil
}
