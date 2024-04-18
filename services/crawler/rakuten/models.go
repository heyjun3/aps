package rakuten

import (
	"crawler/product"
)

const (
	siteCode = "rakuten"
)

func NewRakutenProduct(
	name, productCode, url, jan, shopCode string, price, point int64) (
	*product.Product, error) {

	calcPrice := price - point
	return product.New(
		siteCode,
		shopCode,
		productCode,
		name,
		jan,
		url,
		calcPrice,
	)
}
