package product

import (
	"context"
	"strings"

	"crawler/scrape"

	"github.com/uptrace/bun"
)

type Repository[T scrape.IProduct] struct {
	siteCode string
}

func NewRepository[T scrape.IProduct](siteCode string) Repository[T] {
	return Repository[T]{
		siteCode: siteCode,
	}
}

func (p Repository[T]) GetProduct(ctx context.Context,
	db *bun.DB, productCode, shopCode string) (T, error) {
	var i interface{}
	product := new(Product)
	err := db.NewSelect().
		Model(product).
		Where("(site_code, shop_code, product_code) IN (?)",
			bun.In([][]string{{p.siteCode, shopCode, productCode}})).
		Scan(ctx, product)
	i = product
	result, _ := i.(T)
	return result, err
}

// ここのcodesに型つけたいな
func (p Repository[T]) GetByProductAndShopCodes(ctx context.Context,
	db *bun.DB, codes ...[]string) (scrape.Products, error) {
	records := make([][]string, 0, len(codes))
	for _, code := range codes {
		record := append(code, p.siteCode)
		records = append(records, record)
	}
	var products []*Product
	err := db.NewSelect().
		Model(&products).
		Where("(product_code, shop_code, site_code) IN (?)",
			bun.In(records)).
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
