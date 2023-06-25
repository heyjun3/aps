package nojima

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

type NojimaProduct struct {
	bun.BaseModel `bun:"table:nojima_products"`
	scrape.Product
}

func NewNojimaProduct(name, productCode, url, jan string, price int64) (*NojimaProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, "nojima", price)
	if err != nil {
		return nil, err
	}
	return &NojimaProduct{
		Product: *p,
	}, err
}

func CreateTable(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewCreateTable().
		Model((*NojimaProduct)(nil)).
		IfNotExists().
		Exec(ctx)

	return err
}
