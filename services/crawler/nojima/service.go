package nojima

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*NojimaProduct] {
	return scrape.NewService(NojimaParser{}, &NojimaProduct{}, []*NojimaProduct{})
}
