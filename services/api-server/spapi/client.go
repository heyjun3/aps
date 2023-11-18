package spapi

import (
	"net/url"

	"api-server/spapi/inventory"
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

func (c SpapiClient) GetPricing(ids []string, idType price.IdType) (*price.GetPricingResponse, error) {
	return price.GetPricing(c.URL, ids, idType)
}

func (c SpapiClient) GetLowestPricing(skus []string) (*price.GetLowestPricingResponse, error) {
	return price.GetLowestPricing(c.URL, skus)
}
