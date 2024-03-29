package rakuten

import (
	"context"
	"strings"

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

	calcPrice := price - point
	p, err := scrape.NewProduct(name, productCode, url, jan, shopCode, calcPrice)
	if err != nil {
		return nil, err
	}
	return &RakutenProduct{
		Product: *p,
		point:   point,
	}, nil
}

type Shop struct {
	bun.BaseModel `bun:"table:shops"`
	ID            string `bun:",pk"`
	SiteName      string
	Name          string
	URL           string
	Interval      string
}

type ShopRepository struct{}

func (r ShopRepository) Save(db *bun.DB, ctx context.Context, shops []*Shop) error {
	_, err := db.NewInsert().
		Model(&shops).
		On("CONFLICT (id) DO UPDATE").
		Set(strings.Join([]string{
			"site_name = EXCLUDED.site_name",
			"name = EXCLUDED.name",
			"url = EXCLUDED.url",
			"interval = EXCLUDED.interval",
		}, ",")).
		Exec(ctx)
	return err
}

func (r ShopRepository) GetAll(db *bun.DB, ctx context.Context) ([]*Shop, error) {
	shops := []*Shop{}
	err := db.NewSelect().Model(&shops).Scan(ctx)
	return shops, err
}

func (r ShopRepository) GetByInterval(db *bun.DB, ctx context.Context, interval Interval) ([]Shop, error) {
	shops := []Shop{}
	err := db.NewSelect().Model(&shops).Where("interval = ?", interval.String()).Scan(ctx)
	return shops, err
}

func (r ShopRepository) DeleteShops(ctx context.Context, db *bun.DB, shops []*Shop) error {
	_, err := db.NewDelete().Model(&shops).WherePK().Exec(ctx)
	return err
}

type Interval int

const (
	daily Interval = iota
	weekly
)

func (i Interval) String() string {
	switch i {
	case daily:
		return "daily"
	case weekly:
		return "weekly"
	default:
		return "unknown"
	}
}
