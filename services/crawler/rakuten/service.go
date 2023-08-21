package rakuten

import (
	"context"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*RakutenProduct] {
	return scrape.NewService(RakutenParser{}, &RakutenProduct{}, []*RakutenProduct{})
}

func RunServices() {
	repo := ShopRepository{}
	db := scrape.CreateDBConnection(config.DBDsn)
	shops, err := repo.GetAll(db, context.Background())
	if err != nil {
		logger.Error("error", err)
		return
	}
	for _, s := range shops {
		logger.Info("run service", "shop", s.Name, "url", s.URL)
		NewScrapeService().StartScrape(s.URL, "rakuten")
	}
}
