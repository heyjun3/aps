package kaago

import (
	"log"

	"crawler/product"
	"crawler/scrape"
)

func NewScrapeService(url string) scrape.Service {
	service := scrape.NewService(KaagoParser{},
		scrape.WithCustomRepository(
			product.NewRepository(),
		))
	req, err := generateRequest(0)
	if err != nil {
		log.Fatalln(err)
	}
	service.EntryReq = req
	return service
}
