package price

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
	Offers Offers  `json:"Offers"`
	SKU    *string `json:"SKU"`
}

type Offers []Offer

type Condition struct {
	MyOffer              *bool
	IsFullfilledByAmazon *bool
	IsBuyBoxWinner       *bool
}

func (o Offers) FilterCondition(cond Condition) Offers {
	offers := make(Offers, 0, len(o))
	for _, offer := range o {
		if cond.MyOffer != nil && offer.MyOffer != *cond.MyOffer {
			continue
		}
		if cond.IsFullfilledByAmazon != nil && offer.IsFulfilledByAmazon != *cond.IsFullfilledByAmazon {
			continue
		}
		if cond.IsBuyBoxWinner != nil && offer.IsBuyBoxWinner != *cond.IsBuyBoxWinner {
			continue
		}
		offers = append(offers, offer)
	}
	return offers
}

func (o Offers) Lowest() *Offer {
	if len(o) == 0 {
		return nil
	}
	return &o[0]
}

type Offer struct {
	Points              Point `json:"Points"`
	Price               Price `json:"ListingPrice"`
	MyOffer             bool  `json:"MyOffer"`
	IsFulfilledByAmazon bool  `json:"IsFulfilledByAmazon"`
	IsBuyBoxWinner      bool  `json:"IsBuyBoxWinner"`
}
type Point struct {
	PointsNumber float64 `json:"PointsNumber"`
}
type Price struct {
	Amount *float64 `json:"Amount"`
}

func GetPricing(URL *url.URL, skus []string) (*GetLowestPricingResponse, error) {
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
