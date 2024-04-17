package hikaritv

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"crawler/config"
	"crawler/product"
	"crawler/scrape"
)

const (
	host   = "shop.hikaritv.net"
	scheme = "https"
)

var logger = config.Logger

type HikaritvParser struct {
	scrape.Parser
}

func (p HikaritvParser) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products product.Products
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

	txt := doc.Find(".nextLink").Text()
	if txt == "" {
		logger.Info("Not found next page URL")
		return products, nil
	}
	nextURL, err := url.Parse(req.URL.String())
	if err != nil {
		logger.Error("url parse error", err)
		return products, nil
	}
	query := nextURL.Query()
	page, err := strconv.Atoi(query.Get("currentPage"))
	if err != nil {
		logger.Error("current page convert error", err)
		return products, nil
	}
	query.Set("currentPage", fmt.Sprint(page+1))
	nextURL.RawQuery = query.Encode()
	nextReq, err := http.NewRequest("GET", nextURL.String(), nil)
	if err != nil {
		logger.Error("faild create request", err)
		return products, nil
	}
	return products, nextReq
}

func (p HikaritvParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("[0-9]{12,13}")
	codes := re.FindAllString(doc.Find(".specTable").Text(), -1)
	if len(codes) > 0 {
		return codes[0], nil
	}
	return "", fmt.Errorf("not found jan code")
}
