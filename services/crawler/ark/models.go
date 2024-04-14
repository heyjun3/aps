package ark

import (
	"crawler/product"
)

const (
	siteCode = "ark"
	shopCode = siteCode
)

func newArkProduct(name, productCode, url, jan string, price int64) (*product.Product, error) {
	return product.New(product.Product{
		SiteCode:    siteCode,
		ShopCode:    shopCode,
		ProductCode: productCode,
		Name:        name,
		URL:         url,
		Price:       price,
		Jan:         &jan,
	})
}
