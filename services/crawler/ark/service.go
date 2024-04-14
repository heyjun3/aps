package ark

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeServiceV2() scrape.Service[*product.Product] {
	return scrape.NewService(
		ArkParser{}, &product.Product{}, []*product.Product{},
		scrape.WithCustomRepository(
			product.NewRepository[*product.Product](siteCode)))
}
