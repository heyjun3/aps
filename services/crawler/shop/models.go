package shop

import (
	"context"
	"strings"

	"crawler/config"
	"github.com/uptrace/bun"
)

var (
	logger = config.Logger
)

type Shop struct {
	bun.BaseModel `bun:"table:shops"`
	ID            string `bun:",pk"`
	SiteName      string
	Name          string
	URL           string
	Interval      string
}

type ShopRepository struct{}

func (r ShopRepository) Save(ctx context.Context, db *bun.DB, shops []*Shop) error {
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

func (r ShopRepository) GetAll(ctx context.Context, db *bun.DB) ([]*Shop, error) {
	shops := []*Shop{}
	err := db.NewSelect().Model(&shops).Scan(ctx)
	return shops, err
}

func (r ShopRepository) GetByInterval(ctx context.Context, db *bun.DB, interval Interval) ([]Shop, error) {
	shops := []Shop{}
	err := db.NewSelect().Model(&shops).Where("interval = ?", interval.String()).Scan(ctx)
	return shops, err
}

func (r ShopRepository) GetBySiteName(ctx context.Context, db *bun.DB, siteName string) ([]Shop, error) {
	shops := []Shop{}
	err := db.NewSelect().
		Model(&shops).
		Where("site_name = ?", siteName).
		Scan(ctx)
	return shops, err
}

func (r ShopRepository) GetBySiteNameAndInterval(ctx context.Context, db *bun.DB, siteName string, interval Interval) ([]Shop, error) {
	shops := []Shop{}
	err := db.NewSelect().
		Model(&shops).
		Where("site_name = ?", siteName).
		Where("interval = ?", interval.String()).
		Scan(ctx)
	return shops, err
}

func (r ShopRepository) DeleteShops(ctx context.Context, db *bun.DB, shops []*Shop) error {
	_, err := db.NewDelete().Model(&shops).WherePK().Exec(ctx)
	return err
}

type Interval int

const (
	Daily Interval = iota
	Weekly
)

func (i Interval) String() string {
	switch i {
	case Daily:
		return "daily"
	case Weekly:
		return "weekly"
	default:
		return "unknown"
	}
}
