package pc4u

import (
	"crawler/product"
)

const (
	siteCode = "pc4u"
	shopCode = siteCode
)

func NewPc4uProduct(
	name, productCode, url, jan string, price int64) (*product.Product, error) {
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
