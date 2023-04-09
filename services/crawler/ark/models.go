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

func CreateTable(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewCreateTable().
		Model((*ArkProduct)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func GetByProductCodes(conn *bun.DB, ctx context.Context,
	codes ...string) (scrape.Products, error) {

	var products []*ArkProduct
	err := conn.NewSelect().
		Model(&products).
		Where("product_code IN (?)", bun.In(codes)).
		Scan(ctx)

	var result scrape.Products
	for _, p := range products {
		result = append(result, p)
	}

	return result, err
}
