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
