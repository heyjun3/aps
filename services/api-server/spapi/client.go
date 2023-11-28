package spapi

import (
	"net/url"

	"api-server/spapi/inventory"
	"api-server/spapi/point"
	"api-server/spapi/price"
)

type SpapiClient struct {
	URL *url.URL
}

func NewSpapiClient(URL string) (*SpapiClient, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &SpapiClient{
		URL: u,
	}, nil
}

func (c SpapiClient) InventorySummaries(nextToken string) (*inventory.SummariesResponse, error) {
	return inventory.Summaries(c.URL, nextToken)
}

func (c SpapiClient) GetPricing(skus []string) (*price.GetLowestPricingResponse, error) {
	return price.GetPricing(c.URL, skus)
}

func (c SpapiClient) UpdatePricing(input price.IUpdatePriceInput) error {
	return price.UpdatePricing(c.URL, input)
}

func (c SpapiClient) UpdatePoints(input point.IUpdatePointInput) error {
	return point.UpdatePoints(c.URL, input)
}
