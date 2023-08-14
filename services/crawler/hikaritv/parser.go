package hikaritv

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"crawler/config"
	"crawler/scrape"
)

const (
	host = "shop.hikaritv.net"
	scheme = "https"
)

var logger = config.Logger


type HikaritvParser struct {
	scrape.Parser
}

func (p HikaritvParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products scrape.Products
	doc.Find(".w50p .inner").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".w50p_item_name").Text()
		path, exist := s.Find("a").Attr("href")
		if !exist {
			logger.Info("Not found url", name)
			return
		}
		var paths []string
		for _, p := range strings.Split(path, "/") {
			if p != "" {
				paths = append(paths, p)
			}
		}
		productCode := paths[len(paths)-1]
		URL := url.URL{Scheme: scheme, Host: host, Path: path}
		price, err := scrape.PullOutNumber(s.Find(".numC").Text())
		if err != nil {
			logger.Info("Not found price", "error", err)
			return
		}
		product, err := NewHikaritvProduct(name, "", productCode, URL.String(), price)
		if err != nil {
			logger.Error("error", err)
			return
		}
		products = append(products, product)
	})

	return products, nil
}
