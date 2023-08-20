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
	name, productCode, url, jan, shopCode string, price, point int64) (*RakutenProduct, error) {

	p, err := scrape.NewProduct(name, productCode, url, jan, shopCode, price)
	if err != nil {
		return nil, err
	}
	return &RakutenProduct{
		Product: *p,
		point:   point,
	}, nil
}

func (r *RakutenProduct) calcPrice() {
	r.Price = int64(float64(r.Price)*0.91) - r.point
}

type Shop struct {
	bun.BaseModel `bun:"table:shops"`
	ID            string `bun:",pk"`
	Name          string
	URL           string
}

func GetShops(db *bun.DB, ctx context.Context) ([]Shop, error) {
	shops := []Shop{}
	err := db.NewSelect().Model(shops).Scan(ctx, shops)
	return shops, err
}

func AddShop(db *bun.DB, ctx context.Context, shop Shop) error {
	_, err := db.NewInsert().Model(shop).Exec(ctx)
	return err
}
