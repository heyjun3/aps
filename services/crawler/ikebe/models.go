package ikebe

import (
	"crawler/product"
)

const (
	siteCode = "ikebe"
	shopCode = siteCode
)

func NewIkebeProduct(name, productCode, url, jan string, price int64) (
	*product.Product, error) {
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
