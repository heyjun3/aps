package hikaritv

import (
	"crawler/product"
	"crawler/scrape"
)

func NewScrapeService() scrape.Service {
	return scrape.NewService(HikaritvParser{}, &product.Product{},
		[]*product.Product{}, scrape.WithCustomRepository(
			product.NewRepository(siteCode)))
}
