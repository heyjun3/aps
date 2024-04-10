package ark

import (
	"crawler/product"
)

func newArkProduct(name, productCode, url, jan string, price int64) (*product.Product, error) {
	return product.New(product.Product{
		SiteCode:    "ark",
		ShopCode:    "ark",
		ProductCode: productCode,
		Name:        name,
		URL:         url,
		Price:       price,
		Jan:         &jan,
	})
}
