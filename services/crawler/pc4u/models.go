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

func CreateTable(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewCreateTable().
		Model((*Pc4uProduct)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func GetByProductCodes(conn *bun.DB, ctx context.Context,
	codes ...string)(scrape.Products, error) {

	var products []*Pc4uProduct
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
