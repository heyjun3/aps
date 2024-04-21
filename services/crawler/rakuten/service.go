package rakuten

import (
	"context"

	"crawler/config"
	"crawler/product"
	"crawler/scrape"
	"crawler/shop"
)

var logger = config.Logger

func NewScrapeService() scrape.Service {
	return scrape.NewService(
		RakutenParser{}, &product.Product{}, []*product.Product{},
		scrape.WithCustomRepository(
			product.NewRepository(),
		),
	)
}

func RunServices() {
	repo := shop.ShopRepository{}
	db := scrape.CreateDBConnection(config.DBDsn)
	shops, err := repo.GetBySiteName(context.Background(), db, "rakuten")
	if err != nil {
		logger.Error("error", err)
		return
	}
	for _, s := range shops {
		logger.Info("run service", "shop", s.Name, "url", s.URL)
		NewScrapeService().StartScrape(s.URL, "rakuten")
	}
}

func RunServicesByDaily() {
	repo := shop.ShopRepository{}
	db := scrape.CreateDBConnection(config.DBDsn)
	shops, err := repo.GetBySiteNameAndInterval(context.Background(), db, "rakuten", shop.Daily)
	if err != nil {
		logger.Error("error", err)
		return
	}
	for _, s := range shops {
		logger.Info("run service", "shop", s.Name, "url", s.URL)
		NewScrapeService().StartScrape(s.URL, "rakuten")
	}
}
