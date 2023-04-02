package ikebe

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewIkebeProduct(name, productCode, url, jan string, price int64) *IkebeProduct {
	return &IkebeProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, "ikebe", price),
	}
}

type IkebeProduct struct {
	bun.BaseModel `bun:"table:ikebe_product"`
	scrape.Product
}

func (p *IkebeProduct) Upsert(conn *bun.DB, ctx context.Context) error {
	return scrape.Upsert(conn, ctx, p)
}

func GetByProductCodes(
	conn *bun.DB,ctx context.Context, codes ...string)(scrape.Products, error) {

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
