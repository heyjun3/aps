package ark

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeServiceV2() scrape.Service {
	return scrape.NewService(
		ArkParser{},
		scrape.WithCustomRepository(
			product.NewRepository()))
}
