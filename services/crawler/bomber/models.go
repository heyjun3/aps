package bomber

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

type BomberProduct struct {
	bun.BaseModel `bun:"table:bomber_products"`
	scrape.Product
}

func NewBomberProduct(name, productCode, url, jan string, price int64) (*BomberProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "bomber", price)
	if err != nil {
		return nil, err
	}
	return &BomberProduct{
		Product: *p,
	}, nil
}
