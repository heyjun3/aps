package scrape

import (
	"context"

	"crawler/product"
	"github.com/uptrace/bun"
)

type ProductRepositoryInterface[T product.IProduct] interface {
	BulkUpsert(context.Context, *bun.DB, product.Products) error
	GetProduct(context.Context, *bun.DB, string, string) (T, error)
	GetByProductAndShopCodes(context.Context, *bun.DB, ...[]string) (product.Products, error)
}

type RunServiceHistoryRepository struct{}

func (r RunServiceHistoryRepository) Save(ctx context.Context, db *bun.DB, history *RunServiceHistory) (*RunServiceHistory, error) {
	_, err := db.NewInsert().
		Model(history).
		On("CONFLICT (id) DO UPDATE").
		Set("shop_name = EXCLUDED.shop_name").
		Set("url = EXCLUDED.url").
		Set("status = EXCLUDED.status").
		Set("started_at = EXCLUDED.started_at").
		Set("ended_at = EXCLUDED.ended_at").
		Exec(ctx, history)
	return history, err
}
