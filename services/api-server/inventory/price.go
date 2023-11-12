package inventory

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

type IPrice interface{}

func CastIPrices[T *CurrentPrice](prices []T) []IPrice {
	iprices := make([]IPrice, 0, len(prices))
	for _, price := range prices {
		iprices = append(iprices, IPrice(price))
	}
	return iprices
}

var _ IPrice = (*Price)(nil)

type Price struct {
	SellerSku *string   `bun:"seller_sku,pk"`
	Amount    *int      `bun:"amount"`
	Point     *int      `bun:"point"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewPrice(sku *string, price, point *int) (*Price, error) {
	if slices.Contains([]interface{}{sku, price, point}, nil) {
		return nil, errors.New("contians nil in args")
	}
	return &Price{
		SellerSku: sku,
		Amount:    price,
		Point:     point,
	}, nil
}

var _ IPrice = (*CurrentPrice)(nil)

type CurrentPrice struct {
	bun.BaseModel `bun:"table:current_prices"`
	Price
}

func NewCurrentPrice(sku *string, price, point *int) (*CurrentPrice, error) {
	p, err := NewPrice(sku, price, point)
	if err != nil {
		return nil, err
	}
	return &CurrentPrice{
		Price: *p,
	}, nil
}

type CurrentPrices []*CurrentPrice

var _ IPrice = (*LowestPrice)(nil)

type LowestPrice struct {
	bun.BaseModel `bun:"table:lowest_prices"`
	Price
}

func NewLowestPrice(sku *string, price, point *int) (*LowestPrice, error) {
	p, err := NewPrice(sku, price, point)
	if err != nil {
		return nil, err
	}
	return &LowestPrice{
		Price: *p,
	}, nil
}

type PriceRepository struct{}

func (r PriceRepository) Save(ctx context.Context, db *bun.DB, prices []IPrice) error {
	_, err := db.NewInsert().
		Model(&prices).
		On("CONFLICT (seller_sku) DO UPDATE").
		Set(strings.Join([]string{
			"amount = EXCLUDED.amount",
			"point = EXCLUDED.point",
			"updated_at = current_timestamp",
		}, ",")).
		Returning("NULL").
		Exec(ctx)
	slog.Error("error", "detail", err)
	return err
}
