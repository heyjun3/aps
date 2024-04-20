package bomber

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service {
	return scrape.NewService(
		BomberParser{}, &product.Product{}, []*product.Product{},
		scrape.WithCustomRepository(
			product.NewRepository(siteCode)))
}
