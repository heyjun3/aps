package product

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:mws_products"`
	Asin          string `bun:"asin,pk"`
	Filename      string `bun:"filename,pk"`
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

func (p ProductRepository) GetFilenames(ctx context.Context) ([]Product, error) {
	var products []Product
	err := p.DB.NewSelect().DistinctOn("filename").Model(&products).Scan(ctx)
	return products, err
}
