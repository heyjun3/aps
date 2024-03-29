package ark

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewArkProduct(name, productCode, url, jan string, price int64) (*ArkProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "ark", price)
	if err != nil {
		return nil, err
	}
	return &ArkProduct{
		Product: *p,
	}, nil
}

type ArkProduct struct {
	bun.BaseModel `bun:"table:ark_products"`
	scrape.Product
}
