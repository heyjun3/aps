package bomber

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
	host   = "www.pc-bomber.co.jp"
)

type BomberParser struct{}

func (p BomberParser) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products product.Products
	doc.Find(".pane-main .block-thumbnail-h--goods.js-enhanced-ecommerce-item").Each(func(i int, s *goquery.Selection) {
		nameElem := s.Find(".js-enhanced-ecommerce-goods-name")
		name := nameElem.Text()
		path, exist := nameElem.Attr("href")
		URL, err := url.Parse(path)
		if !exist || err != nil {
			logger.Info("not found url", "name", name, "error", err)
			return
		}
		URL.Scheme = scheme
		URL.Host = host

		var paths []string
		for _, v := range strings.Split(path, "/") {
			if v != "" {
				paths = append(paths, v)
			}
		}
		productCode := paths[len(paths)-1]
		price, err := scrape.PullOutNumber(s.Find(".block-thumbnail-h--price").Text())
		if err != nil {
			logger.Warn("filed parse price", "error", err)
			return
		}

		product, err := NewBomberProduct(name, productCode, URL.String(), "", price)
		if err != nil {
			logger.Warn("filed create bomber product", "error", err)
			return
		}
		products = append(products, product)
	})

	nextReq, err := p.parseNextPageURL(doc)
	if err != nil {
		logger.Error("failed parse next page url", err)
		return products, nil
	}

	return products, nextReq
}

func (p BomberParser) parseNextPageURL(doc *goquery.Document) (*http.Request, error) {
	path, exist := doc.Find(".pager-next a").Attr("href")
	if !exist {
		return nil, fmt.Errorf("not found next page url")
	}
	URL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	URL.Scheme = scheme
	URL.Host = host

	return http.NewRequest("GET", URL.String(), nil)
}

func (p BomberParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("[0-9]{12,13}")

	codes := re.FindAllString(doc.Find(".detail_goods_name2").Text(), -1)
	if len(codes) > 0 {
		return codes[0], nil
	}

	return "", fmt.Errorf("not found jan code")
}
