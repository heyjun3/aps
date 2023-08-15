package hikaritv

import "crawler/scrape"

func NewScrapeService() scrape.Service[*HikaritvProduct] {
	return scrape.NewService(HikaritvParser{}, &HikaritvProduct{}, []*HikaritvProduct{})
}
