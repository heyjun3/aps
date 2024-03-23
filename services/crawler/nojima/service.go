package nojima

import (
	"fmt"
	"net/url"
	"time"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(opts ...scrape.Option[*NojimaProduct]) scrape.Service[*NojimaProduct] {
	return scrape.NewService(NojimaParser{}, &NojimaProduct{}, []*NojimaProduct{}, opts...)
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
	service := NewScrapeService(scrape.WithFileId[*NojimaProduct](fileId))
	for _, u := range urls {
		service.StartScrape(u, "nojima")
	}
}
