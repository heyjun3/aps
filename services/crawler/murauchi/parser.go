package murauchi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"crawler/product"
	"crawler/scrape"
)

const (
	host   = "www.murauchi.com"
	scheme = "https"
	path   = "MCJ-front-web/WH/front/Default.do"
)

type MurauchiParser struct{}

func (p MurauchiParser) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	reader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products product.Products
	doc.Find(".window_item").Each(func(i int, s *goquery.Selection) {
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

	isLastPage := true
	doc.Find(".search_paging a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "最後") {
			isLastPage = false
		}
	})
	if isLastPage {
		logger.Info("next page is nothing")
		return products, nil
	}

	nextRequest, err := generateRequestFromPreviousRequest(req)
	if err != nil {
		logger.Error("error", err)
		return products, nil
	}

	return products, nextRequest
}

func (p MurauchiParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	return doc.Find("form[name=mail_form] input[name=jan_code]").AttrOr("value", ""), nil
}

func (p MurauchiParser) FindCategories(r io.ReadCloser) ([]string, error) {
	reader := transform.NewReader(r, japanese.EUCJP.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	categories := []string{}
	doc.Find("#categories section a").Each(func(i int, s *goquery.Selection) {
		link, exist := s.Attr("href")
		if !exist {
			logger.Info("href tag is not found")
			return
		}
		paths := []string{}
		for _, path := range strings.Split(link, "/") {
			if path != "" {
				paths = append(paths, path)
			}
		}
		category := paths[2]
		categories = append(categories, category)
	})
	return categories, nil
}

func generateRequestFromPreviousRequest(pre *http.Request) (*http.Request, error) {
	category := pre.Header.Get("x-category")
	page := pre.Header.Get("x-current")
	if category == "" || page == "" {
		return nil, fmt.Errorf("category not found error")
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	return generateRequest(p+1, category)
}

func generateRequest(page int, category string) (*http.Request, error) {
	keyword, err := scrape.Utf8ToSjis("　")
	if err != nil {
		return nil, err
	}
	values := map[string]string{
		"mode":         "graphic",
		"pageNumber":   fmt.Sprint(page),
		"searchType":   "keyword",
		"sortOrder":    "1",
		"categoryNo":   category,
		"type":         "COMMODITY_LIST",
		"keyword":      keyword,
		"listCount":    "120",
		"handlingType": "0",
	}
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}
	body := strings.NewReader(form.Encode())

	u := url.URL{Scheme: scheme, Host: host, Path: path}
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-category", category)
	req.Header.Set("x-current", fmt.Sprint(page))

	return req, nil
}
