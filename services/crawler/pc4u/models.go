package pc4u

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewPc4uProduct(name, productCode, url, jan string, price int64) *Pc4uProduct {
	return &Pc4uProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, "pc4u", price),
	}
}

type Pc4uProduct struct {
	bun.BaseModel `bun:"table:pc4u_products"`
	scrape.Product
}
