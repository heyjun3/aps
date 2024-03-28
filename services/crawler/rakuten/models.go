package rakuten

import (
	"github.com/uptrace/bun"

	"crawler/scrape"
)

type RakutenProduct struct {
	bun.BaseModel `bun:"table:rakuten_products"`
	scrape.Product
	point int64
}

func NewRakutenProduct(
	name, productCode, url, jan, shopCode string, price, point int64) (*RakutenProduct, error) {

	calcPrice := price - point
	p, err := scrape.NewProduct(name, productCode, url, jan, shopCode, calcPrice)
	if err != nil {
		return nil, err
	}
	return &RakutenProduct{
		Product: *p,
		point:   point,
	}, nil
}
