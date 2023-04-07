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
