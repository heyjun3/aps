package ark

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*ArkProduct] {
	return scrape.NewService(ArkParser{}, &ArkProduct{}, []*ArkProduct{})
}
