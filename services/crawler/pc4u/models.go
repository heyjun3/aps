package pc4u

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewPc4uProduct(name, productCode, url, jan string, price int64) *Pc4uProduct {
	return &Pc4uProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, "pc4u", price),
	}
}

type Pc4uProduct struct {
	bun.BaseModel `bun:"table:pc4u_products"`
	scrape.Product
}

func (p *Pc4uProduct) Upsert(conn *bun.DB, ctx context.Context) error {
	return scrape.Upsert(conn, ctx, p)
}

func GetByProductCodes(conn *bun.DB,
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
