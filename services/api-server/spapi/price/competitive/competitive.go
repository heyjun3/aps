package competitive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetCompetitivePricingResponse struct {
	response
}

type response struct {
	Payload []Payload `json:"payload"`
}
type Payload struct {
	Asin    string  `json:"ASIN"`
	Product Product `json:"Product"`
	Status  string  `json:"status"`
}
type Product struct {
	CompetitivePricing CompetitivePricing `json:"CompetitivePricing"`
	SalesRankings      []SalesRank        `json:"SalesRankings"`
}
type CompetitivePricing struct {
	CompetitivePrices []CompetitivePrice `json:"CompetitivePrices"`
}
type CompetitivePrice struct {
	Price Price `json:"Price"`
}
type Price struct {
	LandedPrice  Amount `json:"LandedPrice"`
	ListingPrice Amount `json:"ListingPrice"`
}
type Amount struct {
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
	fmt.Println(string(body))
	var getCompetitivePricingResponse GetCompetitivePricingResponse
	if err := json.Unmarshal(body, &getCompetitivePricingResponse); err != nil {
		return nil, err
	}
	return &getCompetitivePricingResponse, nil
}
