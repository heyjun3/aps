package murauchi

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"crawler/scrape"
)

const (
	host = "www.murauchi.com"
	scheme = "https"
)

type MurauchiParser struct{}

func (p MurauchiParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products scrape.Products
	doc.Find(".window_item").Each(func(i int, s*goquery.Selection) {
		nameElem := s.Find("h2 a")
		name := nameElem.Text()
		href, exist := nameElem.Attr("href")
		URL, err := url.Parse(href)
		if !exist || err != nil {
			logger.Info("Not Found url")
			return
		}
		var paths []string
		for _, path := range strings.Split(URL.Path, "/") {
			if path != "" {
				paths = append(paths, path)
			}
		}
		productId := paths[2]
		URL = &url.URL{Host: host, Scheme: scheme, Path: strings.Join(paths[:3], "/")}
		priceText := s.Find(".price_p .price").Text()
		price, err := scrape.PullOutNumber(priceText)
		if err != nil {
			logger.Info("Not Found price")
			return
		}

		sold := s.Find(".stock span").Text()
		if sold := strings.TrimSpace(sold); sold == "販売終了" || sold == "予約中" {
			logger.Info("product is sold out")
			return
		}
		product, err := NewMurauchiProduct(name, productId, URL.String(), "", price)
		if err != nil {
			logger.Info("error", err)
			return
		}
		products = append(products, product)
	})

	return products, &http.Request{}
}
