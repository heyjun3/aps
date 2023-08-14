package hikaritv

import (
	"crawler/scrape"
)

type HikaritvProduct struct {
	scrape.Product
}

func NewHikaritvProduct(name, jan, productCode, url string, price int64) (*HikaritvProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "hikaritv", price)
	if err != nil {
		return nil, err
	}
	return &HikaritvProduct{
		Product: *p,
	}, nil
}
