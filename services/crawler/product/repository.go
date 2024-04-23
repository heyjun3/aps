package product

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct{}

func NewRepository() Repository {
	return Repository{}
}

func (p Repository) GetProduct(ctx context.Context,
	db *bun.DB, siteCode, shopCode, productCode string) (*Product, error) {
	product := new(Product)
	err := db.NewSelect().
		Model(product).
		Where("(site_code, shop_code, product_code) IN (?)",
			bun.In([][]string{{siteCode, shopCode, productCode}})).
		Scan(ctx, product)
	return product, err
}

func (p Repository) GetByCodes(ctx context.Context,
	db *bun.DB, codes []Code) (Products, error) {
	records := make([][]string, 0, len(codes))
	for _, code := range codes {
		record := []string{code.SiteCode, code.ShopCode, code.ProductCode}
		records = append(records, record)
	}
	var products []*Product
	err := db.NewSelect().
		Model(&products).
		Where("(site_code, shop_code, product_code) IN (?)",
			bun.In(records)).
		Order("product_code ASC").
		Scan(ctx, &products)
	return ConvToProducts(products), err
}

func (p Repository) BulkUpsert(ctx context.Context, db *bun.DB,
	ps Products) error {
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
