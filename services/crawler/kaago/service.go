package kaago

import (
	"log"

	"crawler/scrape"
)

func NewScrapeService(url string) scrape.Service[*KaagoProduct] {
	service := scrape.NewService(KaagoParser{}, &KaagoProduct{}, []*KaagoProduct{})
	req, err := generateRequest(0)
	if err != nil {
		log.Fatalln(err)
	}
	service.EntryReq = req
	return service
}
