package product

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:mws_products"`
	Asin          string    `bun:"asin,pk" json:"asin"`
	Filename      string    `bun:"filename,pk" json:"filename"`
	Title         *string   `json:"title"`
	Jan           *string   `json:"jan"`
	Unit          *int64    `json:"unit"`
	Price         *int64    `json:"price"`
	Cost          *int64    `json:"cost"`
	FeeRate       *float64  `json:"fee_rate"`
	ShippingFee   *int64    `bun:"shipping_fee" json:"shipping_fee"`
	Profit        *int64    `json:"profit"`
	ProfitRate    *float64  `bun:"profit_rate" json:"profit_rate"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	URL           *string   `bun:"url" json:"url"`
}

type RenderChart struct {
	Chart ChartData `bun:"render_data,type:jsonb"`
	Drops int       `bun:"sales_drops_90"`
}

type ProductWithChart struct {
	Product
	RenderChart
}

type ProductRepository struct {
	DB *bun.DB
}

type Condition struct {
	Unit       *int64
	Profit     *int64
	ProfitRate *float64
}

func NewCondition(profit, unit int64, profit_rate float64) *Condition {
	return &Condition{
		Profit:     &profit,
		ProfitRate: &profit_rate,
		Unit:       &unit,
	}
}

func (p ProductRepository) Save(ctx context.Context, products []Product) error {
	_, err := p.DB.NewInsert().Model(&products).Exec(ctx)
	return err
}

func (p ProductRepository) Upserts(ctx context.Context, products []Product) error {
	_, err := p.DB.NewInsert().
		Model(&products).
		ExcludeColumn("profit", "profit_rate").
		On("CONFLICT (asin, filename) DO UPDATE").
		Set(`
			title = EXCLUDED.title,
			jan = EXCLUDED.jan,
			unit = EXCLUDED.unit,
			price = EXCLUDED.price,
			cost = EXCLUDED.cost,
			fee_rate = EXCLUDED.fee_rate,
			shipping_fee = EXCLUDED.shipping_fee,
			url = EXCLUDED.url
		`).
		Exec(ctx)
	return err
}

func (p ProductRepository) GetProduct(ctx context.Context) ([]Product, error) {
	var products []Product
	err := p.DB.NewSelect().Model(&products).Scan(ctx)
	return products, err
}

func (p ProductRepository) GetCounts(ctx context.Context) (map[string]int, error) {
	var total, price, fee int
	err := p.DB.NewSelect().
		Model((*Product)(nil)).
		ColumnExpr("count(*)").
		ColumnExpr("count(?)", bun.Ident("price")).
		ColumnExpr("count(?)", bun.Ident("fee_rate")).
		Scan(ctx, &total, &price, &fee)
	return map[string]int{"total": total, "price": price, "fee": fee}, err
}

func (p ProductRepository) GetFilenames(ctx context.Context) ([]string, error) {
	var filenames []string
	subquery := p.DB.NewSelect().
		Model((*Product)(nil)).
		Column("filename").
		DistinctOn("filename").
		Where("profit IS NULL")

	err := p.DB.NewSelect().
		Model((*Product)(nil)).
		Column("filename").
		DistinctOn("filename").
		Where("filename NOT IN (?)", subquery).
		Order("filename ASC").
		Scan(ctx, &filenames)
	return filenames, err
}

type searchCondition struct {
	filename                string
	page                    int
	limit                   int
	minProfit               int
	minProfitRate           float32
	maxUnit                 int
	minQuarterlySalesVolume int
}

func NewSearchCondition(filename string, opts ...searchConditionOption) searchCondition {
	c := searchCondition{
		filename:                filename,
		page:                    1,
		limit:                   100,
		minProfit:               200,
		minProfitRate:           0.1,
		maxUnit:                 10,
		minQuarterlySalesVolume: 3,
	}
	for _, opt := range opts {
		c = opt(c)
	}
	return c
}

type searchConditionOption func(searchCondition) searchCondition

func SearchConditionWithPage(page int) searchConditionOption {
	return func(c searchCondition) searchCondition {
		c.page = page
		return c
	}
}
func SearchConditionWithLimit(limit int) searchConditionOption {
	return func(c searchCondition) searchCondition {
		c.limit = limit
		return c
	}
}

func (p ProductRepository) GetProductWithChartBySearchCondition(
	ctx context.Context, c searchCondition) ([]ProductWithChart, int, error) {
	if c.page < 1 {
		return nil, 0, fmt.Errorf("page is over 1")
	}
	offset := (c.page - 1) * c.limit
	var product []ProductWithChart
	count, err := p.DB.NewSelect().
		ColumnExpr("p.*").
		ColumnExpr("k.render_data").
		ColumnExpr("k.sales_drops_90").
		TableExpr("mws_products AS p").
		Join("JOIN keepa_products AS k ON k.asin = p.asin").
		Where("p.filename = ?", c.filename).
		Where("p.profit >= ?", c.minProfit).
		Where("p.profit_rate >= ?", c.minProfitRate).
		Where("p.unit <= ?", c.maxUnit).
		Where("k.sales_drops_90 > ?", c.minQuarterlySalesVolume).
		Where("k.render_data IS NOT NULL").
		OrderExpr("p.profit DESC").
		Limit(c.limit).
		Offset(offset).
		ScanAndCount(ctx, &product)
	return product, count, err
}

func (p ProductRepository) DeleteIfCondition(ctx context.Context, condition *Condition) error {
	_, err := p.DB.NewDelete().
		Model((*Product)(nil)).
		WhereOr("profit < ?", condition.Profit).
		WhereOr("profit_rate < ?", condition.ProfitRate).
		WhereOr("unit > ?", condition.Unit).
		Exec(ctx)
	return err
}

func (p ProductRepository) DeleteIfConditionWithKeepa(ctx context.Context) error {
	_, err := p.DB.NewDelete().
		Model((*Product)(nil)).
		TableExpr("keepa_products as keepa").
		Where("product.asin = keepa.asin").
		Where("keepa.sales_drops_90 <= ?", 3).
		Exec(ctx)
	return err
}

func (p ProductRepository) DeleteByFilename(ctx context.Context, filename string) error {
	_, err := p.DB.NewDelete().Model((*Product)(nil)).Where("filename = ?", filename).Exec(ctx)
	return err
}

func (p ProductRepository) RefreshGeneratedColumns(ctx context.Context) error {
	_, err := p.DB.NewUpdate().
		Model((*Product)(nil)).
		Set("price = price").
		WhereOr("profit IS NULL").
		WhereOr("profit_rate IS NULL").
		Exec(ctx)
	return err
}
