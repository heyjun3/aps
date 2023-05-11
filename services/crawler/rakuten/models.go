package rakuten

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

type RakutenProduct struct {
	bun.BaseModel `bun:"table:rakuten_products"`
	scrape.Product
	point int64
}

func NewRakutenProduct(
	name, productCode, url, jan, shopCode string, price, point int64) *RakutenProduct {
	return &RakutenProduct{
		Product: *scrape.NewProduct(name, productCode, url, jan, shopCode, price),
		point:   point,
	}
}

func (r *RakutenProduct) calcPrice() {
	r.Price = int64(float64(r.Price)*0.91) - r.point
}

func CreateTable(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewCreateTable().
		Model((*RakutenProduct)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func getByProductCodes(shopCode string) func(*bun.DB, context.Context, ...string) (scrape.Products, error) {
	return func(conn *bun.DB, ctx context.Context, codes ...string) (scrape.Products, error) {

		var products []*RakutenProduct
		err := conn.NewSelect().
			Model(&products).
			Where("product_code IN (?)", bun.In(codes)).
			Where("shop_code = ?", shopCode).
			Order("product_code ASC").
			Scan(ctx)

		var result scrape.Products
		for _, p := range products {
			result = append(result, p)
		}
		return result, err
	}
}
