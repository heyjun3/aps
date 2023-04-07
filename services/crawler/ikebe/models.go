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

func CreateTable(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewCreateTable().
		Model((*IkebeProduct)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}
