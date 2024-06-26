package ark

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"crawler/product"
	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host   = "www.ark-pc.co.jp"
)

type ArkParser struct {
	scrape.Parser
}

func (p ArkParser) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	return p.ConvToReq(p.ProductList(r, req.URL.String()))
}

func (p ArkParser) ProductList(r io.ReadCloser, u string) (product.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}
	var products product.Products
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
		var paths []string
		for _, v := range splitPath {
			if v != "" {
				paths = append(paths, v)
			}
		}
		productId := paths[len(paths)-1]

		price, err := scrape.PullOutNumber(s.Find(".itemprice .price").Text())
		if err != nil {
			logger.Info("Not Found price")
			return
		}
		coupon, _ := scrape.PullOutNumber(s.Find(".price_diff_2.auto_coupon").Text())
		discountedPrice := price - coupon

		product, err := newArkProduct(name, productId, URL.String(), "", discountedPrice)
		if err != nil {
			logger.Error("error", err)
			return
		}
		products = append(products, product)
	})

	path, exist := doc.Find("#listnavi_next a[href]").Attr("href")
	nextURL, err := url.Parse(path)
	if !exist || err != nil {
		logger.Info("Not Found next page URL")
		return products, ""
	}
	nextURL.Scheme = scheme
	nextURL.Host = host

	return products, nextURL.String()
}

func (p ArkParser) Product(res io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("[0-9]{12,13}")
	jancodes := re.FindAllString(doc.Find(".jancode").Text(), -1)

	if len(jancodes) > 0 {
		return jancodes[0], nil
	}
	return "", fmt.Errorf("not found jan code")
}
