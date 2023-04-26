package rakuten

import (
	"crawler/scrape"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host   = "item.rakuten.co.jp"
)

type RakutenParser struct{}

func (p RakutenParser) ProductList(r io.ReadCloser) (scrape.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}

	var products scrape.Products
	doc.Find(".dui-card.searchresultitem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".title-link--3Ho6z").Text()
		if name == "" {
			logger.Info("Not Found product name")
			return
		}

		path, exist := s.Find(".image a").Attr("href")
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
		productId := paths[len(paths)-1]
		shopId := paths[len(paths)-2]

		price, err := scrape.PullOutNumber(s.Find(".price--OX_YW").Text())
		if err != nil {
			logger.Info("Not Found price", "name", name, "url", URL.String())
			return
		}
		re := regexp.MustCompile("[0-9,]+ポイント")
		point, err := scrape.PullOutNumber(re.FindString(s.Find(".points--AHzKn span").Text()))
		if err != nil {
			logger.Info("Not Found point", "name", name, "url", URL.String())
		}

		products = append(products,
			NewRakutenProduct(name, productId, URL.String(), "", shopId, price, point))
	})

	nextURL, err := p.scrapeNextURL(doc, products[len(products)-1].GetPrice())
	if err != nil {
		logger.Error("not found next page url %s", err)
		return products, ""
	}
	return products, nextURL
}

func (p RakutenParser) scrapeNextURL(doc *goquery.Document, minPrice int64) (string, error) {
	nextURL, exist := doc.Find(".item.-next.nextPage").Attr("href")
	if exist {
		URL, err := url.Parse(nextURL)
		if err != nil {
			return "", err
		}
		return URL.String(), nil
	}
	currentPage := doc.Find(".item.-active.currentPage").Text()
	if currentPage != "150" {
		logger.Info("curretn page is not 150")
		return "", nil
	}
	currentURL, exist := doc.Find("link[rel=canonical]").Attr("href")
	if !exist {
		return "", errors.New("not found current url")
	}
	URL, err := url.Parse(currentURL)
	if err != nil {
		return "", err
	}
	query := URL.Query()
	query.Set("max", fmt.Sprint(minPrice))
	query.Set("p", "1")
	query.Set("s", "12")
	URL.RawQuery = query.Encode()

	return URL.String(), nil
}

func (p RakutenParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	jan, exist := doc.Find("meta[itemprop=gtin13]").Attr("content")
	if !exist {
		return "", errors.New("not found jan code")
	}
	re := regexp.MustCompile("^[0-9]{12,13}$")
	match := re.Match([]byte(jan))
	if match {
		return jan, nil
	}
	return "", fmt.Errorf("not match jan code value: %s", jan)
}
