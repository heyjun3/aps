package ark

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewArkProduct(name, productCode, url, jan string, price int64) *ArkProduct {
	return &ArkProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, "ark", price),
	}
}

type ArkProduct struct {
	bun.BaseModel `bun:"table:ark_products"`
	scrape.Product
}

func (p *ArkProduct) Upsert(conn *bun.DB, ctx context.Context) error {
	return scrape.Upsert(conn, ctx, p)
}

func GetByProductCodes(conn *bun.DB,
	ctx context.Context, codes ...string) (scrape.Products, error) {

	var arkProducts []ArkProduct
	err := conn.NewSelect().
		Model(&arkProducts).
		Where("product_code IN (?)", bun.In(codes)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var products scrape.Products
	for i := 0; i < len(arkProducts); i++ {
		products = append(products, &arkProducts[i])
	}
	return products, nil
}
