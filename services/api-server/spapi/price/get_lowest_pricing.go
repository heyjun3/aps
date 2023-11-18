package price

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetLowestPricingResponse struct{}

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
	fmt.Println(string(body))

	return nil, nil
}
