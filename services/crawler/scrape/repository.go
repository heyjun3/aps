package scrape

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository interface {
	GetByProductCodes(conn *bun.DB, ctx context.Context, codes ...string) (Products, error)
	// Upsert(conn *bun.DB, ctx context.Context, p Product) error
}

type ProductRepository struct {}

// func (r ProductRepository) Upsert(conn *bun.DB, ctx context.Context, p Product) error {
// 	_, err := conn.NewInsert().
// 		Model(p).
// 		On("CONFLICT (shop_code, product_code) DO UPDATE").
// 		Set(`
// 		name = EXCLUDED.name,
// 		jan = EXCLUDED.jan,
// 		price = EXCLUDED.price,
// 		url = EXCLUDED.url
// 		`).
// 		Exec(ctx)
// 	return err
// }
