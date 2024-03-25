package pc4u

import (
	"strings"
	"time"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(opts ...scrape.Option[*Pc4uProduct]) scrape.Service[*Pc4uProduct] {
	return scrape.NewService(Pc4uParser{}, &Pc4uProduct{}, []*Pc4uProduct{}, opts...)
}

func ScrapeAll(shop string) {
	urls := []string{
		"https://www.pc4u.co.jp/view/category/outlet",
		"https://www.pc4u.co.jp/view/search",
	}
	fileId := strings.Join([]string{shop, scrape.TimeToStr(time.Now())}, "_")
	service := NewScrapeService(scrape.WithFileId[*Pc4uProduct](fileId))
	for _, url := range urls {
		service.StartScrape(url, shop)
	}
}
