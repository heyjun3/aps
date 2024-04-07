package product

import (
	"context"
	"strings"

	"crawler/scrape"

	"github.com/uptrace/bun"
)

type Repository[T scrape.IProduct] struct{}

func (p Repository[T]) GetProduct(ctx context.Context,
	db *bun.DB, productCode, shopCode string) (T, error) {
	var i interface{}
	product := new(Product)
	err := db.NewSelect().
		Model(product).
		Where("product_code = ?", productCode).
		Where("shop_code = ?", shopCode).
		Scan(ctx, product)
	i = product
	result, _ := i.(T)
	return result, err
}

func (p Repository[T]) GetByProductAndShopCodes(ctx context.Context,
	db *bun.DB, codes ...[]string) (scrape.Products, error) {
	var products []*Product
	err := db.NewSelect().
		Model(&products).
		Where("(product_code, shop_code) IN (?)", bun.In(codes)).
		Order("product_code ASC").
		Scan(ctx, &products)
	return scrape.ConvToProducts(products), err
}

func (p Repository[T]) BulkUpsert(ctx context.Context, db *bun.DB,
	ps scrape.Products) error {
	mapProduct := map[string]scrape.IProduct{}
	for _, v := range ps {
		mapProduct[v.GetProductCode()] = v
	}
	var products scrape.Products
	for _, v := range mapProduct {
		products = append(products, v)
	}

	_, err := db.NewInsert().
		Model(&products).
		On("CONFLICT (site_code, shop_code, product_code) DO UPDATE").
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
