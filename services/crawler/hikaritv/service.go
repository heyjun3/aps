package hikaritv

import (
	"crawler/product"
	"crawler/scrape"
)

func NewScrapeService() scrape.Service[*product.Product] {
	return scrape.NewService(HikaritvParser{}, &product.Product{},
		[]*product.Product{}, scrape.WithCustomRepository(
			product.NewRepository[*product.Product](siteCode)))
}
