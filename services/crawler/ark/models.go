package ark

import (
	"crawler/product"
)

const (
	siteCode = "ark"
	shopCode = siteCode
)

func newArkProduct(name, productCode, url, jan string, price int64) (*product.Product, error) {
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
