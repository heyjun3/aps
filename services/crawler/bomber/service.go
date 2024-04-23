package bomber

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service {
	return scrape.NewService(
		BomberParser{},
		scrape.WithCustomRepository(
			product.NewRepository()))
}
