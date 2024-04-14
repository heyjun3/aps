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
		product.Product{
			SiteCode:    siteCode,
			ShopCode:    shopCode,
			ProductCode: productCode,
			Name:        name,
			Jan:         &jan,
			Price:       price,
			URL:         url,
		},
	)
}
