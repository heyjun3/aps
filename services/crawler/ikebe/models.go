package ikebe

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewIkebeProduct(name, productCode, url, jan string, price int64) *IkebeProduct {
	return &IkebeProduct{
		BaseProduct: *scrape.NewProduct(name, productCode, url, jan, "ikebe", price),
	}
}

type IkebeProduct struct {
	bun.BaseModel `bun:"table:ikebe_product"`
	scrape.BaseProduct
}
