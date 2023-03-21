package ikebe

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

// type IkebeProduct struct {
// 	models.IkebeProduct
// }

// func NewIkebeProduct(name, productCode, url, jan string, price int64) *IkebeProduct {
// 	isJan := true
// 	if jan == "" {
// 		isJan = false
// 	}
// 	return &IkebeProduct{
// 		models.IkebeProduct{
// 			Name:        null.StringFrom(name),
// 			Jan:         null.NewString(jan, isJan),
// 			Price:       null.Int64From(price),
// 			ShopCode:    "ikebe",
// 			ProductCode: productCode,
// 			URL:         null.StringFrom(url),
// 		},
// 	}
// }

// func (p *IkebeProduct) GenerateMessage(filename string) ([]byte, error) {
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

// func (p *IkebeProduct) GetProductCode() string {
// 	return p.ProductCode
// }

// func (p *IkebeProduct) GetJan() string {
// 	return p.Jan.String
// }

// func (p *IkebeProduct) GetURL() string {
// 	return p.URL.String
// }

// func (p *IkebeProduct) IsValidJan() bool {
// 	return p.Jan.Valid
// }

// func (p *IkebeProduct) SetJan(jan string) {
// 	p.Jan = null.StringFrom(jan)
// }

type IkebeProductRepository struct{
	scrape.ProductRepository
}

func (r IkebeProductRepository) GetByProductCodes(conn *bun.DB,
	ctx context.Context, codes ...string) (scrape.Products, error) {

	var ikebeProducts []IkebeProduct
	err := conn.NewSelect().
		Model(&ikebeProducts).
		Where("product_code IN (?)", bun.In(codes)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var products scrape.Products
	for i := 0; i < len(ikebeProducts); i++ {
		products = append(products, &ikebeProducts[i])
	}
	return products, nil
}

// func (repo IkebeProductRepository) bulkUpsert(conn *sql.DB, products ...*IkebeProduct) error {
// 	strs := []string{}
// 	args := []interface{}{}
// 	for i, p := range products {
// 		d := i * 6
// 		str := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", d+1, d+2, d+3, d+4, d+5, d+6)
// 		strs = append(strs, str)

// 		args = append(args, p.Name.String)
// 		args = append(args, p.Jan.String)
// 		args = append(args, p.Price.Int64)
// 		args = append(args, p.ShopCode)
// 		args = append(args, p.ProductCode)
// 		args = append(args, p.URL.String)
// 	}
// 	stmt := fmt.Sprintf(`INSERT INTO ikebe_product (name, jan, price, shop_code, product_code, url) 
// 						VALUES %s ON CONFLICT (shop_code, product_code) DO UPDATE SET 
// 						name = excluded.name, jan = excluded.jan, price = excluded.price, 
// 						url = excluded.url RETURNING url, name;`, strings.Join(strs, ","))
// 	_, err := conn.Exec(stmt, args...)
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

