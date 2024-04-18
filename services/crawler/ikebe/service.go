package ikebe

import (
	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

const (
	scheme = "https"
	host   = "www.ikebe-gakki.com"
)

func NewScrapeService() scrape.Service {
	return scrape.NewService(IkebeParser{}, &product.Product{},
		[]*product.Product{}, scrape.WithCustomRepository(
			product.NewRepository(siteCode),
		))
}
