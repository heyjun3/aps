package product

import (
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Chart struct {
	Date  string  `json:"date" bun:"type:date"`
	Rank  float64 `json:"rank"`
	Price float64 `json:"price"`
}

type ChartData struct {
	Data []Chart `json:"data"`
}

type Keepa struct {
	bun.BaseModel `bun:"keepa_products"`
	Asin          string             `bun:"asin,pk"`
	Drops         int                `bun:"sales_drops_90"`
	Prices        map[string]float64 `bun:"price_data,type:jsonb"`
	Ranks         map[string]float64 `bun:"rank_data,type:jsonb"`
	Charts        ChartData          `bun:"render_data,type:jsonb"`
	Created       time.Time          `bun:",type:date,nullzero,notnull,default:current_timestamp"`
	Modified      time.Time          `bun:",type:date,nullzero,notnull,default:current_timestamp"`
}
type Keepas []*Keepa

func (k Keepas) Asins() []string {
	asins := make([]string, len(k))
	for i, keepa := range k {
		asins[i] = keepa.Asin
	}
	return asins
}

type KeepaRepository struct {
	DB *bun.DB
}

func (k KeepaRepository) Save(ctx context.Context, Keepas []*Keepa) error {
	_, err := k.DB.
		NewInsert().
		Model(&Keepas).
		On("CONFLICT (asin) DO UPDATE").
		Set(strings.Join([]string{
			"sales_drops_90 = EXCLUDED.sales_drops_90",
			"price_data = EXCLUDED.price_data",
			"rank_data =  EXCLUDED.rank_data",
			"render_data = EXCLUDED.render_data",
			"modified = now()",
		}, ",")).
		Exec(ctx)
	return err
}

func (k KeepaRepository) Get(ctx context.Context) (*Keepa, error) {
	keepa := new(Keepa)
	err := k.DB.NewSelect().Model(keepa).Limit(1).Scan(ctx)
	return keepa, err
}

func (k KeepaRepository) GetByAsins(ctx context.Context, asins []string) (Keepas, error) {
	keepas := make([]*Keepa, 0, len(asins))
	err := k.DB.NewSelect().
		Model(&keepas).
		Where("asin IN (?)", bun.In(asins)).
		Order("asin").
		Scan(ctx)
	return keepas, err
}

func (k KeepaRepository) GetCounts(ctx context.Context) (map[string]int, error) {
	now := time.Now().Format("2006-01-02")

	var total, modified int
	err := k.DB.NewSelect().
		Model((*Keepa)(nil)).
		ColumnExpr("count(*)").
		ColumnExpr("count(? = ? or NULL)", bun.Ident("modified"), now).
		Scan(ctx, &total, &modified)

	return map[string]int{"total": total, "modified": modified}, err
}
