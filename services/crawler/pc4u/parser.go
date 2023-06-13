package pc4u

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host   = "www.pc4u.co.jp"
)

type Pc4uParser struct {
	scrape.Parser
}

func (p Pc4uParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	return p.ConvToReq(p.ProductList(r, req.URL.String()))
}

func (p Pc4uParser) ProductList(r io.ReadCloser, u string) (scrape.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}

	isSold := false
	var products scrape.Products
	doc.Find(".big-item-list__item").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".big-item-list__name a").Text()
		if name == "" {
			logger.Info("Not Found product name", "value", name)
			return
		}

		path, exist := s.Find(".big-item-list__name a").Attr("href")
		URL, err := url.Parse(path)
		if !exist || err != nil {
			logger.Info("Not Found url", "value", path)
			return
		}
		URL.Scheme = scheme
		URL.Host = host

		paths := strings.Split(URL.Path, "/")
		productId := paths[len(paths)-1]

		priceText := s.Find(".big-item-list__price").Text()
		price, err := scrape.PullOutNumber(priceText)
		if err != nil {
			logger.Info("Not Found price", "value", priceText)
			return
		}

		sold := s.Find(".big-item-list__soldout").Text()
		if sold != "" {
			logger.Info("product is sold out")
			isSold = true
			return
		}
		products = append(products, NewPc4uProduct(name, productId, URL.String(), "", price))
	})

	if isSold {
		logger.Info("products contain sold out product")
		return products, ""
	}

	nextURL, err := p.nextPageURL(doc)
	if err != nil {
		logger.Error("not found next page url", err)
		return products, ""
	}

	return products, nextURL
}

func (p Pc4uParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`[0-9]{12,13}`)
	itemDescription := doc.Find(".item-description__content").Text()
	janCodes := re.FindAllString(itemDescription, -1)
	if len(janCodes) > 0 {
		return janCodes[0], nil
	}

	return "", fmt.Errorf("not found jan")
}

func (p Pc4uParser) nextPageURL(doc *goquery.Document) (string, error) {
	var nextPath = ""
	var current = false
	doc.Find(".pager__item a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		path, exist := s.Attr("href")
		if !exist {
			current = true
			return true
		}
		if current {
			nextPath = path
			return false
		}
		return true
	})

	nextURL, err := url.Parse(nextPath)
	if err != nil || nextPath == "" {
		return "", fmt.Errorf("not found next page URL: %s", nextPath)
	}
	nextURL.Scheme = scheme
	nextURL.Host = host

	return nextURL.String(), nil
}
