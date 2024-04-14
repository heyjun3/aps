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
	return product.New(product.Product{
		SiteCode:    siteCode,
		ShopCode:    shopCode,
		ProductCode: productCode,
		Name:        name,
		Jan:         &jan,
		Price:       price,
		URL:         url,
	})
}
