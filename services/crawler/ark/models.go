package ark

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewArkProduct(name, productCode, url, jan string, price int64) *ArkProduct {
	return &ArkProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, "ark", price),
	}
}

type ArkProduct struct {
	bun.BaseModel `bun:"table:ark_products"`
	scrape.Product
}
