package bomber

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*BomberProduct] {
	return scrape.NewService(BomberParser{}, &BomberProduct{}, []*BomberProduct{})
}
