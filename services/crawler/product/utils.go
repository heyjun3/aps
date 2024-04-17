package product

func NewTestProduct(
	name, productCode, url, jan, shopCode string, price int64) *Product {
	p, _ := New(
		Product{
			SiteCode:    "testSite",
			ShopCode:    shopCode,
			ProductCode: productCode,
			Name:        name,
			URL:         url,
			Jan:         &jan,
			Price:       price,
		},
	)
	return p
}
