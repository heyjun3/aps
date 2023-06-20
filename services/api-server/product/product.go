package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(sqldb, pgdialect.New())
}

type Product struct {
	bun.BaseModel `bun:"table:mws_products"`
	Asin          string `bun:"asin,pk"`
	Filename      string `bun:"filename,pk" json:"filename"`
	Title         string
	Jan           string
	Unit          int64
	Price         int64
	Cost          int64
	FeeRate       float64
	ShippingFee   int64 `bun:"shipping_fee"`
	Profit        int64
	ProfitRate    float64   `bun:"profit_rate"`
	CreatedAt     time.Time `bun:"created_at"`
	URL           string    `bun:"url"`
}

type ProductRepository struct {
	DB *bun.DB
}

func (p ProductRepository) Save(ctx context.Context, products []Product) (error) {
	_, err := p.DB.NewInsert().Model(&products).Exec(ctx)
	return err
}

func (p ProductRepository) GetFilenames(ctx context.Context) ([]string, error) {
	var filenames []string
	err := p.DB.NewSelect().Model((*Product)(nil)).Column("filename").DistinctOn("filename").Order("filename ASC").Scan(ctx, &filenames)
	return filenames, err
}

func (p ProductRepository) DeleteByFilename(ctx context.Context, filename string) error {
	_, err := p.DB.NewDelete().Model((*Product)(nil)).Where("filename = ?", filename).Exec(ctx)
	return err
}
