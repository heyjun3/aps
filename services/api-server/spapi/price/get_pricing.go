package price

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/slog"
)

type GetPricingResponse struct {
	Payload []Payload `json:"payload"`
}

type Payload struct {
	Status    string            `json:"status"`
	SellerSKU string            `json:"SellerSKU"`
	Product   Product `json:"Product"`
}

type Product struct {
	Offers []Offers `json:"Offers"`
}

type Offers struct {
	BuyingPrice BuyingPrice `json:"BuyingPrice"`
}

type BuyingPrice struct {
	ListingPrice Price  `json:"ListingPrice"`
	Points       Points `json:"Points"`
}

type Price struct {
	CurrencyCode string  `json:"CurrencyCode"`
	Amount       float64 `json:"Amount"`
}

type Points struct {
	PointsNumber int64 `json:"PointsNumber"`
}
type IdType int

const (
	Asin IdType = iota
	Sku
)

func (t IdType) String() string {
	switch t {
	case Asin:
		return "Asin"
	case Sku:
		return "Sku"
	default:
		return "unknown"
	}
}

func GetPricing(URL *url.URL, ids []string, idType IdType) (*GetPricingResponse, error) {
	query := url.Values{}
	query.Set("ids", strings.Join(ids, ","))
	query.Set("id_type", idType.String())
	URL.Path = "get-pricing"
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
	var getPricingResponse GetPricingResponse
	if err := json.Unmarshal(body, &getPricingResponse); err != nil {
		slog.Error("err", err)
		return nil, err
	}
	return &getPricingResponse, nil
}
