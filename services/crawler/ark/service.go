package ark

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*ArkProduct] {
	return scrape.NewService(ArkParser{}, &ArkProduct{}, []*ArkProduct{})
}

func NewScrapeServiceV2() scrape.Service[*product.Product] {
	return scrape.NewService(
		ArkParser{}, &product.Product{}, []*product.Product{},
		scrape.WithCustomRepository(product.Repository[*product.Product]{}))
}
