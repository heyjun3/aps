package nojima

import (
	"fmt"
	"net/url"
	"time"

	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(
	opts ...scrape.Option[*product.Product]) scrape.Service {
	return scrape.NewService(
		NojimaParser{}, &product.Product{}, []*product.Product{}, opts...)
}

func ScrapeAll() {
	var urls []string
	endpoint := "https://online.nojima.co.jp/app/catalog/list/init"
	URL, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	for i := 100; i < 120; i++ {
		q := URL.Query()
		q.Set("searchCategoryCode", fmt.Sprint(i))
		q.Set("immediateDeliveryDispFlg", "1")
		URL.RawQuery = q.Encode()
		urls = append(urls, URL.String())
	}
	fileId := "nojima_" + scrape.TimeToStr(time.Now())
	service := NewScrapeService(
		scrape.WithFileId[*product.Product](fileId),
		scrape.WithCustomRepository(
			product.NewRepository(),
		),
	)
	for _, u := range urls {
		service.StartScrape(u, "nojima")
	}
}
