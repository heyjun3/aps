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

func (p *Pc4uProduct) Upsert(conn *bun.DB, ctx context.Context) error {
	return scrape.Upsert(conn, ctx, p)
}
