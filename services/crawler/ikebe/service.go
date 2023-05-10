package ikebe

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

const (
	scheme = "https"
	host   = "www.ikebe-gakki.com"
)

func NewScrapeService() *scrape.Service {
	return &scrape.Service{
		FetchProductByProductCodes: scrape.GetByProductCodes([]*IkebeProduct{}),
		Parser:                     IkebeParser{},
	}
}
