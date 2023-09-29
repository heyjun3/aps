package scrape

import (
	"context"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
)

type ProductRepository[T IProduct] struct {
	product  T
	products []T
}

func NewProductRepository[T IProduct](p T, ps []T) ProductRepository[T] {
	return ProductRepository[T]{
		product:  p,
		products: ps,
	}
}

func (p ProductRepository[T]) GetProduct(ctx context.Context,
	db *bun.DB, productCode, shopCode string) (T, error) {

	product := reflect.New(reflect.ValueOf(p.product).Elem().Type()).Interface().(T)
	err := db.NewSelect().
		Model(product).
		Where("product_code = ?", productCode).
		Where("shop_code = ?", shopCode).
		Scan(ctx, product)

	return product, err
}

func (p ProductRepository[T]) GetByProductAndShopCodes(ctx context.Context, db *bun.DB, codes ...[]string) (Products, error) {
	products := p.products

	err := db.NewSelect().
		Model(&products).
		Where("(product_code, shop_code) IN (?)", bun.In(codes)).
		Order("product_code ASC").
		Scan(ctx, &products)

	return ConvToProducts(products), err
}

func (p ProductRepository[T]) BulkUpsert(ctx context.Context, db *bun.DB, ps Products) error {
	mapProduct := map[string]IProduct{}
	for _, v := range ps {
		mapProduct[v.GetProductCode()] = v
	}
	var products Products
	for _, v := range mapProduct {
		products = append(products, v)
	}

	_, err := db.NewInsert().
		Model(&products).
		On("CONFLICT (shop_code, product_code) DO UPDATE").
		Set(strings.Join([]string{
			"name = EXCLUDED.name",
			"jan = EXCLUDED.jan",
			"price = EXCLUDED.price",
			"url = EXCLUDED.url",
		}, ",")).
		Returning("NULL").
		Exec(ctx)

	return err
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
