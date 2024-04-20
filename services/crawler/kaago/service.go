package kaago

import (
	"log"

	"crawler/product"
	"crawler/scrape"
)

func NewScrapeService(url string) scrape.Service {
	service := scrape.NewService(KaagoParser{},
		&product.Product{}, []*product.Product{}, scrape.WithCustomRepository(
			product.NewRepository(siteCode),
		))
	req, err := generateRequest(0)
	if err != nil {
		log.Fatalln(err)
	}
	service.EntryReq = req
	return service
}
