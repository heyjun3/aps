package nojima

import (
	"fmt"
	"net/url"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*NojimaProduct] {
	return scrape.NewService(NojimaParser{}, &NojimaProduct{}, []*NojimaProduct{})
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
	service := NewScrapeService()
	for _, u := range urls {
		service.StartScrape(u, "nojima")
	}
}
