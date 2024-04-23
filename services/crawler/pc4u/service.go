package pc4u

import (
	"context"
	"strings"
	"time"

	"crawler/config"
	"crawler/product"
	"crawler/scrape"
	"crawler/shop"
)

var logger = config.Logger

func NewScrapeService(opts ...scrape.Option) scrape.Service {
	return scrape.NewService(
		Pc4uParser{}, opts...)
}

func ScrapeAll(shopName string) {
	shopRepository := shop.ShopRepository{}
	db := scrape.CreateDBConnection(config.DBDsn)
	shops, err := shopRepository.GetBySiteName(context.Background(), db, "pc4u")
	if err != nil {
		panic(err)
	}
	fileId := strings.Join([]string{shopName, scrape.TimeToStr(time.Now())}, "_")
	service := NewScrapeService(
		scrape.WithFileId(fileId),
		scrape.WithCustomRepository(
			product.NewRepository(),
		),
	)
	for _, s := range shops {
		service.StartScrape(s.URL, shopName)
	}
}
