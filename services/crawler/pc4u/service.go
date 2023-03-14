package pc4u

import (
	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() *scrape.Service {
	return &scrape.Service{
		Repo:   Pc4uProductRepository{},
		Parser: Pc4uParser{},
	}
}
