package ark

import (
	"io"
	"net/url"
	"strings"

	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host = "www.ark-pc.co.jp"
)

type ArkParser struct {}

func (p ArkParser) ProductList(r io.ReadCloser) (scrape.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}
	var products scrape.Products
	doc.Find(".item_listbox").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".itemname1 a").Text()
		if name == "" {
			logger.Info("Not Found product name")
			return
		}
		path, exist := s.Find(".itemname1 a").Attr("href")
		URL, err := url.Parse(path)
		if !exist || err != nil {
			logger.Info("Not Found URL")
			return
		}
		URL.Scheme = scheme
		URL.Host = host

		splitPath := strings.Split(path, "/")
		productId := splitPath[len(splitPath)-1]

		price, err := scrape.PullOutNumber(s.Find(".itemprice .price").Text())
		if err != nil {
			logger.Info("Not Found price")
			return
		}
		coupon, err := scrape.PullOutNumber(s.Find(".price_diff_2_coupon").Text())
		if err != nil {
			logger.Info("Not Found coupon")
			return
		}
		discountedPrice := price - coupon
		

	})

}