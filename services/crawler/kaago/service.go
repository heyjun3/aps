package kaago

import ( 
	"crawler/scrape"
)

func NewScrapeService() scrape.Service[*KaagoProduct]{
	return scrape.NewService(KaagoParser{}, &KaagoProduct{}, []*KaagoProduct{})
}
