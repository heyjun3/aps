package murauchi

import (
	"crawler/product"
)

const (
	siteCode = "murauchi"
	shopCode = siteCode
)

func NewMurauchiProduct(
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
