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
		product.Product{
			SiteCode:    siteCode,
			ShopCode:    shopCode,
			ProductCode: productCode,
			Name:        name,
			URL:         url,
			Jan:         &jan,
			Price:       price,
		},
	)
}
