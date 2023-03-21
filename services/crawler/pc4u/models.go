package pc4u

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewPc4uProduct(name, productCode, url, jan string, price int64) *Pc4uProduct {
	return &Pc4uProduct{
		BaseProduct: *scrape.NewProduct(name, productCode, url, jan, "pc4u", price),
	}
}

type Pc4uProduct struct {
	bun.BaseModel `bun:"table:pc4u_products"`
	scrape.BaseProduct
}

// type Pc4uProduct struct {
	// models.Pc4uProduct
// }

// func NewPc4uProduct(name, productCode, url, jan string, price int64) *Pc4uProduct {
// 	isJan := true
// 	if jan == "" {
// 		isJan = false
// 	}
// 	return &Pc4uProduct{
// 		models.Pc4uProduct{
// 			Name:        null.StringFrom(name),
// 			Jan:         null.NewString(jan, isJan),
// 			Price:       null.Int64From(price),
// 			ShopCode:    "pc4u",
// 			ProductCode: productCode,
// 			URL:         null.StringFrom(url),
// 		},
// 	}
// }

// func (p *Pc4uProduct) GenerateMessage(filename string) ([]byte, error) {
// 	if !p.Jan.Valid {
// 		return nil, fmt.Errorf("jan code isn't valid %s", p.ProductCode)
// 	}
// 	if !p.Price.Valid {
// 		return nil, fmt.Errorf("price isn't valid %s", p.ProductCode)
// 	}
// 	if !p.URL.Valid {
// 		return nil, fmt.Errorf("url isn't valid %s", p.ProductCode)
// 	}
// 	m := scrape.NewMWSSchema(filename, p.Jan.String, p.URL.String, p.Price.Int64)
// 	message, err := json.Marshal(m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return message, err
// }

// func (p *Pc4uProduct) GetProductCode() string {
// 	return p.ProductCode
// }

// func (p *Pc4uProduct) GetJan() string {
// 	return p.Jan.String
// }

// func (p *Pc4uProduct) GetURL() string {
// 	return p.URL.String
// }

// func (p *Pc4uProduct) IsValidJan() bool {
// 	return p.Jan.Valid
// }

// func (p *Pc4uProduct) SetJan(jan string) {
// 	p.Jan = null.StringFrom(jan)
// }

type Pc4uProductRepository struct{
	scrape.ProductRepository
}

func (r Pc4uProductRepository) GetByProductCodes(conn *bun.DB,
	ctx context.Context, codes ...string) (scrape.Products, error) {

	var pc4uProducts []Pc4uProduct
	err := conn.NewSelect().
		Model(&pc4uProducts).
		Where("product_code IN (?)", bun.In(codes)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var products scrape.Products
	for i := 0; i < len(pc4uProducts); i++ {
		products = append(products, &pc4uProducts[i])
	}
	return products, nil
}
