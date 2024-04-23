package hikaritv

import (
	"crawler/product"
	"crawler/scrape"
)

func NewScrapeService() scrape.Service {
	return scrape.NewService(HikaritvParser{}, scrape.WithCustomRepository(
		product.NewRepository()))
}
