package murauchi

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

type MurauchiProduct struct {
	bun.BaseModel `bun:"table:murauchi_products"`
	scrape.Product
}

func NewMurauchiProduct(name, productCode, url, jan string, price int64) (*MurauchiProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "murauchi", price)
	if err != nil {
		return nil, err
	}
	return &MurauchiProduct{
		Product: *p,
	}, nil
}
