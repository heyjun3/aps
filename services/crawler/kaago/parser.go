package kaago

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"crawler/config"
	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = "https"
	host   = "kaago.com"
	path   = "ajax/catalog/list/init"
)

var logger = config.Logger

type KaagoParser struct {
	previousPage int
}

func (p KaagoParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	var resp KaagoResp
	var products scrape.Products

	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		logger.Error("response decode error", err)
		return products, nil
	}

	if p.previousPage == int(resp.CurrentPage) {
		return products, nil
	}
	p.previousPage = int(resp.CurrentPage)
	logger.Info(fmt.Sprintf("current page: %d", p.previousPage))

	for _, p := range resp.ProductList {
		if err := ValidateKaagoRespProduct(p); err != nil {
			logger.Info("product contains zero value field", "err", err)
			continue
		}
		u, err := url.Parse(p.URL)
		if err != nil {
			logger.Error("url parse error", err, "value", p.URL)
			continue
		}
		u.Scheme = scheme
		u.Host = host
		products = append(products, NewKaagoProduct(p.Name, p.ProductCode, u.String(), p.ProductCode, p.ShopCode, p.Price))
	}

	nextReq, err := p.generateRequest(resp.CurrentPage)
	if err != nil {
		return products, nil
	}

	return products, nextReq
}

func (p KaagoParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	return doc.Find("#commodityCode").AttrOr("value", ""), nil
}

func (p KaagoParser) generateRequest(currentPage int64) (*http.Request, error) {
	values := map[string]string{
		"categorycode": "0",
		"currentPage":  fmt.Sprint(currentPage),
		"hasStock":     "1",
		"shopcode":     "",
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
	return req, nil
}
