package nojima

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"crawler/product"
	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host   = "online.nojima.co.jp"
)

type NojimaParser struct {
	scrape.Parser
}

func (p NojimaParser) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	return p.ConvToReq(p.ProductList(r, req.URL.String()))
}

func (p NojimaParser) ProductList(r io.ReadCloser, requestURL string) (product.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("reponse parse error", err)
		return nil, ""
	}

	var products product.Products
	doc.Find(".shouhinlist").Each(func(i int, s *goquery.Selection) {
		name := strings.Join(strings.Fields(s.Find(".textOverflowShohinmei a[href]").Text()), "")
		if name == "" {
			logger.Info("Not Found product name")
			return
		}

		path, exist := s.Find(".textOverflowShohinmei a[href]").Attr("href")
		URL, err := url.Parse(path)
		if !exist || err != nil {
			logger.Info("Not Found url", "name", name)
			return
		}
		URL.Scheme = scheme
		URL.Host = host

		var paths []string
		for _, p := range strings.Split(URL.Path, "/") {
			if p != "" {
				paths = append(paths, p)
			}
		}
		productCode := paths[len(paths)-1]
		jan := productCode

		price, err := scrape.PullOutNumber(s.Find(".price").Text())
		if err != nil {
			logger.Info("Not Found price", "name", "url", URL.String())
			return
		}
		product, err := NewNojimaProduct(name, productCode, URL.String(), jan, price)
		if err != nil {
			logger.Error("error", err)
			return
		}
		products = append(products, product)
	})

	isLastPage := true
	doc.Find(".listwaku").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "次へ" {
			isLastPage = false
		}
	})
	if isLastPage {
		return products, ""
	}

	nextURL, err := p.generateNextURL(doc, requestURL)
	if err != nil {
		return products, ""
	}

	return products, nextURL
}

func (p NojimaParser) generateNextURL(doc *goquery.Document, requestURL string) (string, error) {
	nextURL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}

	page, err := scrape.PullOutNumber(doc.Find(".listwakuselect").First().Text())
	if err != nil {
		return "", err
	}

	query := nextURL.Query()
	query.Set("currentPage", fmt.Sprint(page+1))
	nextURL.RawQuery = query.Encode()

	return nextURL.String(), nil
}

func (p NojimaParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	u, exist := doc.Find("link[rel=canonical]").Attr("href")
	if !exist {
		return "", fmt.Errorf("not found canonical url")
	}
	var path []string
	for _, v := range strings.Split(u, "/") {
		if v != "" {
			path = append(path, v)
		}
	}
	return path[len(path)-1], nil
}
