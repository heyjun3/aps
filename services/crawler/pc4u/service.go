package pc4u

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*Pc4uProduct] {
	return scrape.NewService(Pc4uParser{}, &Pc4uProduct{}, []*Pc4uProduct{})
}
