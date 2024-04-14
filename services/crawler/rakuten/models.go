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
		product.Product{
			SiteCode:    siteCode,
			ShopCode:    shopCode,
			ProductCode: productCode,
			Name:        name,
			URL:         url,
			Jan:         &jan,
			Price:       calcPrice,
		},
	)
}
