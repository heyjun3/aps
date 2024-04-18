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

func NewScrapeService(opts ...scrape.Option[*product.Product]) scrape.Service {
	return scrape.NewService(
		Pc4uParser{}, &product.Product{}, []*product.Product{}, opts...)
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
		scrape.WithFileId[*product.Product](fileId),
		scrape.WithCustomRepository(
			product.NewRepository(siteCode),
		),
	)
	for _, s := range shops {
		service.StartScrape(s.URL, shopName)
	}
}
