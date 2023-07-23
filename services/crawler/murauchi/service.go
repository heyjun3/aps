package murauchi

import (
	"log"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(category string) scrape.Service[*MurauchiProduct] {
	service := scrape.NewService(MurauchiParser{}, &MurauchiProduct{}, []*MurauchiProduct{})
	req, err := generateRequest(0, category)
	if err != nil {
		log.Fatalln(err)
	}
	service.EntryReq = req
	return service
}
