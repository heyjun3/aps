package hikaritv

import (
	"crawler/product"
)

const (
	siteCode = "hikaritv"
	shopCode = siteCode
)

func NewHikaritvProduct(name, jan, productCode, url string, price int64) (
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
