package spapi

import (
	"net/url"

	"api-server/spapi/inventory"
	"api-server/spapi/price"
	"api-server/spapi/price/lowest"
	"api-server/spapi/price/update"
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

func (c SpapiClient) GetLowestPricing(skus []string) (*lowest.GetLowestPricingResponse, error) {
	return lowest.GetLowestPricing(c.URL, skus)
}

func (c SpapiClient) UpdatePricing(sku string, price int) error {
	return update.Pricing(c.URL, sku, price)
}
