package kaago

import (
	"log"
	"net/http"
	"strings"

	"crawler/scrape"
)

func NewScrapeService(url, payload string) scrape.Service[*KaagoProduct] {
	service := scrape.NewService(KaagoParser{}, &KaagoProduct{}, []*KaagoProduct{})
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	service.EntryReq = req
	return service
}
