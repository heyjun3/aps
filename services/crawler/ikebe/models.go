package ikebe

import (
	"crawler/product"
)

const (
	siteCode = "ikebe"
	shopCode = siteCode
)

func NewIkebeProduct(name, productCode, url, jan string, price int64) (
	*product.Product, error) {
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
	// p, err := scrape.NewProduct(name, productCode, url, jan, "ikebe", price)
	// if err != nil {
	// 	return nil, err
	// }
	// return &IkebeProduct{
	// 	Product: *p,
	// }, nil
}

// type IkebeProduct struct {
// 	bun.BaseModel `bun:"table:ikebe_product"`
// 	scrape.Product
// }
