package ikebe

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

func NewIkebeProduct(name, productCode, url, jan string, price int64) *IkebeProduct {
	return &IkebeProduct{
		BaseProduct: *scrape.NewProduct(name, productCode, url, jan, "ikebe", price),
	}
}

type IkebeProduct struct {
	bun.BaseModel `bun:"table:ikebe_product"`
	scrape.BaseProduct
}

func (p *IkebeProduct) Upsert(conn *bun.DB, ctx context.Context) error {
	return scrape.Upsert(conn, ctx, p)
}