package inventory

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"

	"api-server/spapi/point"
	"api-server/spapi/price"
)

type IPrice interface {
	GetPrice() *int
}

var _ IPrice = (*Price)(nil)

type Price struct {
	SellerSku    *string   `bun:"seller_sku,pk"`
	Amount       *int      `bun:"amount,notnull"`
	Point        *int      `bun:"point,notnull"`
	PercentPoint *int      `bun:"percent_point,notnull"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewPrice(sku *string, price, point *float64) (*Price, error) {
	if slices.Contains([]interface{}{sku, price}, nil) {
		return nil, errors.New("contians nil in args")
	}
	po := func() float64 {
		if point == nil {
			return 0
		}
		return float64(*point)
	}()
	percent := int(math.Round(po / *price * 100))
	return &Price{
		SellerSku:    sku,
		Amount:       Ptr(int(*price)),
		Point:        Ptr(int(po)),
		PercentPoint: &percent,
	}, nil
}

func NewPriceWithPercentPoint(sku *string, price, percentPoint *float64) (*Price, error) {
	if percentPoint == nil {
		return nil, errors.New("percent point is nil. expect int argument")
	}
	p, err := NewPrice(sku, price, nil)
	if err != nil {
		return nil, err
	}
	point := int(math.Round(*price / 100 * *percentPoint))
	p.PercentPoint = Ptr(int(*percentPoint))
	p.Point = &point
	return p, nil
}

func (p *Price) GetPrice() *int {
	return p.Amount
}

var _ IPrice = (*CurrentPrice)(nil)

type CurrentPrices []*CurrentPrice
type CurrentPrice struct {
	bun.BaseModel `bun:"table:current_prices"`
	Price
}

func NewCurrentPrice(sku string, price, point *float64) (*CurrentPrice, error) {
	p, err := NewPrice(&sku, price, point)
	if err != nil {
		return nil, err
	}
	return &CurrentPrice{
		Price: *p,
	}, nil
}

var _ IPrice = (*LowestPrice)(nil)

type LowestPrices []*LowestPrice
type LowestPrice struct {
	bun.BaseModel `bun:"table:lowest_prices"`
	Price
}

func NewLowestPrice(sku string, price, point *float64) (*LowestPrice, error) {
	p, err := NewPrice(&sku, price, point)
	if err != nil {
		return nil, err
	}
	return &LowestPrice{
		Price: *p,
	}, nil
}

type DesiredPrices []*DesiredPrice
type DesiredPrice struct {
	bun.BaseModel `bun:"table:desired_prices"`
	Price
}

func NewDesiredPrice(sku *string, price, percentPoint *float64, lowestPrice LowestPrice) (*DesiredPrice, error) {
	p, err := NewPriceWithPercentPoint(sku, price, percentPoint)
	if err != nil {
		return nil, err
	}

	const DESIRED_PRICE_RATE float64 = 0.9
	if *price < float64(*lowestPrice.Amount)*DESIRED_PRICE_RATE {
		return nil, errors.New("desired price greater than lowest price")
	}

	if *p.SellerSku != *lowestPrice.SellerSku {
		return nil, errors.New("not equals skus")
	}

	return &DesiredPrice{
		Price: *p,
	}, nil
}

func (p DesiredPrice) UpdatePrice() price.UpdatePriceInput {
	return price.UpdatePriceInput{
		Sku:   *p.SellerSku,
		Price: *p.Amount,
	}
}

func (p DesiredPrices) UpdatePoints() []point.UpdatePointInput {
	inputs := make([]point.UpdatePointInput, 0, len(p))
	for _, price := range p {
		inputs = append(inputs, point.UpdatePointInput{
			Sku:          *price.SellerSku,
			PercentPoint: *price.PercentPoint,
		})
	}
	return inputs
}

type PriceRepository[T IPrice] struct{}

func (r PriceRepository[T]) Save(ctx context.Context, db *bun.DB, prices []T) error {
	if len(prices) == 0 {
		return nil
	}
	_, err := db.NewInsert().
		Model(&prices).
		On("CONFLICT (seller_sku) DO UPDATE").
		Set(strings.Join([]string{
			"amount = EXCLUDED.amount",
			"point = EXCLUDED.point",
			"percent_point = EXCLUDED.percent_point",
			"updated_at = current_timestamp",
		}, ",")).
		Returning("NULL").
		Exec(ctx)
	return err
}
