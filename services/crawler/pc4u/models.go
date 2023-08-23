package pc4u

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewPc4uProduct(name, productCode, url, jan string, price int64) (*Pc4uProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "pc4u", price)
	if err != nil {
		return nil, err
	}
	return &Pc4uProduct{
		Product: *p,
	}, nil
}

type Pc4uProduct struct {
	bun.BaseModel `bun:"table:pc4u_products"`
	scrape.Product
}
