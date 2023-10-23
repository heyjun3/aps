package inventory

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

type Inventory struct {
	bun.BaseModel   `bun:"table:inventories"`
	Asin            string `json:"asin" bun:"asin"`
	FnSku           string `json:"fnSku" bun:"fnsku"`
	SellerSku       string `json:"sellerSku" bun:"seller_sku,pk"`
	Condition       string `json:"condition" bun:"condition"`
	LastUpdatedTime string `json:"lastUpdatedTime" bun:"-"`
	ProductName     string `json:"productName" bun:"product_name"`
	TotalQuantity   int    `json:"totalQuantity" bun:"quantity"`
	Price           *int   `bun:"price"`
	Point           *int   `bun:"point"`
	LowestPrice     *int   `bun:"lowest_price"`
	LowestPoint     *int   `bun:"lowest_point"`
}

type InventoryRepository struct{}

func (r InventoryRepository) Save(ctx context.Context, db *bun.DB, inventories []*Inventory) error {
	_, err := db.NewInsert().
		Model(&inventories).
		On("CONFLICT (seller_sku) DO UPDATE").
		Set(strings.Join([]string{
			"asin = EXCLUDED.asin",
			"fnsku = EXCLUDED.fnsku",
			"condition = EXCLUDED.condition",
			"product_name = EXCLUDED.product_name",
			"quantity = EXCLUDED.quantity",
			"price = EXCLUDED.price",
			"point = EXCLUDED.point",
			"lowest_price = EXCLUDED.lowest_price",
			"lowest_point = EXCLUDED.lowest_point",
		}, ",")).
		Exec(ctx)
	return err
}
