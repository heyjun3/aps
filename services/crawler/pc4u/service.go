package pc4u

import (
	"context"
	"strings"
	"time"

	"crawler/config"
	"crawler/scrape"
	"crawler/shop"
)

var logger = config.Logger

func NewScrapeService(opts ...scrape.Option[*Pc4uProduct]) scrape.Service[*Pc4uProduct] {
	return scrape.NewService(Pc4uParser{}, &Pc4uProduct{}, []*Pc4uProduct{}, opts...)
}

func ScrapeAll(shopName string) {
	shopRepository := shop.ShopRepository{}
	db := scrape.CreateDBConnection(config.DBDsn)
	shops, err := shopRepository.GetBySiteName(context.Background(), db, "pc4u")
	if err != nil {
		panic(err)
	}
	fileId := strings.Join([]string{shopName, scrape.TimeToStr(time.Now())}, "_")
	service := NewScrapeService(scrape.WithFileId[*Pc4uProduct](fileId))
	for _, s := range shops {
		service.StartScrape(s.URL, shopName)
	}
}
