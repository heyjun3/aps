package competitive

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetCompetitivePricingResponse struct {
	response
}

func (r GetCompetitivePricingResponse) LandedPrices() []*LandedProduct {
	landedPrices := make([]*LandedProduct, 0, len(r.Payload))
	for _, p := range r.Payload {
		landedPrices = append(landedPrices, p.landedProduct())
	}
	return landedPrices
}

type response struct {
	Payload []Payload `json:"payload"`
}
type Payload struct {
	Asin    string  `json:"ASIN"`
	Product Product `json:"Product"`
	Status  string  `json:"status"`
}
type LandedProduct struct {
	Asin          string
	Status        string
	LandedPrice   *Price
	ListingPrice  *Price
	SalesRankings []SalesRank
}

func (p Payload) landedProduct() *LandedProduct {
	rankings := p.Product.SalesRankings
	prices := p.Product.CompetitivePricing.CompetitivePrices
	if len(prices) == 0 {
		return &LandedProduct{
			Asin:          p.Asin,
			Status:        p.Status,
			LandedPrice:   nil,
			ListingPrice:  nil,
			SalesRankings: rankings,
		}
	}
	return &LandedProduct{
		Asin:          p.Asin,
		Status:        p.Status,
		LandedPrice:   prices[0].Prices.LandedPrice,
		ListingPrice:  prices[0].Prices.ListingPrice,
		SalesRankings: rankings,
	}
}

type Product struct {
	CompetitivePricing CompetitivePricing `json:"CompetitivePricing"`
	SalesRankings      []SalesRank        `json:"SalesRankings"`
}
type CompetitivePricing struct {
	CompetitivePrices []CompetitivePrice `json:"CompetitivePrices"`
}
type CompetitivePrice struct {
	Prices Prices `json:"Price"`
}
type Prices struct {
	LandedPrice  *Price `json:"LandedPrice"`
	ListingPrice *Price `json:"ListingPrice"`
}
type Price struct {
	CurrencyCode string `json:"CurrencyCode"`
	Amount       int    `json:"Amount"`
}
type SalesRank struct {
	ProductCategoryId string `json:"ProductCategoryId"`
	Rank              int    `json:"Rank"`
}

func GetCompetitivePricing(URL *url.URL, asins []string) (*GetCompetitivePricingResponse, error) {
	if len(asins) == 0 {
		return nil, errors.New("asins must contain at least on letter")
	}
	if len(asins) > 20 {
		return nil, errors.New("expect asins less than 20 length")
	}

	query := url.Values{}
	query.Set("asins", strings.Join(asins, ","))
	URL.Path = "competitive-pricing"
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
	var getCompetitivePricingResponse GetCompetitivePricingResponse
	if err := json.Unmarshal(body, &getCompetitivePricingResponse); err != nil {
		return nil, err
	}
	return &getCompetitivePricingResponse, nil
}
