package bomber

import (
	"crawler/product"
)

func NewBomberProduct(name, productCode, url, jan string,
	price int64) (*product.Product, error) {
	return product.New(product.Product{
		SiteCode:    "bomber",
		ShopCode:    "bomber",
		ProductCode: productCode,
		Name:        name,
		Jan:         &jan,
		Price:       price,
		URL:         url,
	})
}
