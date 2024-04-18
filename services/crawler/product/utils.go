package product

func NewTestProduct(
	name, productCode, url, jan, shopCode string, price int64) *Product {
	p, _ := New(
		"testSite",
		shopCode,
		productCode,
		name,
		jan,
		url,
		price,
	)
	return p
}
