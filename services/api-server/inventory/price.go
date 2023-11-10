package inventory

import (
	"time"

	"github.com/uptrace/bun"
)

type Price struct {
	SellerSku string    `bun:"seller_sku,pk"`
	Amount    int       `bun:"amount"`
	Point     int       `bun:"point"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type CurrentPrice struct {
	bun.BaseModel `bun:"table:prices"`
	Price
}

type LowestPrice struct {
	bun.BaseModel `bun:"table:lowest_prices"`
	Price
}
