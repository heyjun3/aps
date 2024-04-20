package bomber

import (
	"crawler/product"
)

const (
	siteCode = "bomber"
	shopCode = siteCode
)

func NewBomberProduct(name, productCode, url, jan string,
	price int64) (*product.Product, error) {
	return product.New(
		siteCode,
		shopCode,
		productCode,
		name,
		jan,
		url,
		price,
	)
}
