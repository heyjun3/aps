package price

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
