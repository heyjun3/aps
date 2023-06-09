package kaago

import (
	"encoding/json"
	"io"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

type KaagoParser struct{}

func (p KaagoParser) ProductList(r io.ReadCloser, url string) (scrape.Products, string) {
	var resp KaagoResp

	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		logger.Error("response decode error", err)
		return nil, ""
	}

	var products scrape.Products
	for _, p := range resp.ProductList {
		if err := ValidateKaagoRespProduct(p); err != nil {
			logger.Info("product contains zero value field", "err", err)
			continue
		}
		products = append(products, NewKaagoProduct(p.Name, p.ProductCode, p.URL, p.ProductCode, p.ShopCode, p.Price))
	}
	print(len(products))

	return products, ""
}

func (p KaagoParser) Product(r io.ReadCloser) (string, error) {
	return "", nil
}
