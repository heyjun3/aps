package ark

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() *scrape.Service {
	return &scrape.Service{
		FetchProductByProductCodes: GetByProductCodes,
		Parser: ArkParser{},
	}
}
