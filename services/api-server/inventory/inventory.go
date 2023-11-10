package inventory

import (
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"

	inventory "api-server/spapi/inventory"
)

type Inventory struct {
	bun.BaseModel `bun:"table:inventories"`
	*inventory.Inventory
	CurrentPrice *CurrentPrice `bun:"rel:has-one,join:seller_sku=seller_sku"`
	LowestPrice  *LowestPrice  `bun:"rel:has-one,join:seller_sku=seller_sku"`
	CreatedAt    time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time     `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewInventory(
	asin, fnSku, sellerSku, condition, productName string, totalQuantity int,
) *Inventory {
	return &Inventory{
		Inventory: &inventory.Inventory{
			Asin:          asin,
			FnSku:         fnSku,
			SellerSku:     sellerSku,
			Condition:     condition,
			ProductName:   productName,
			TotalQuantity: totalQuantity,
		},
	}
}

type Inventories []*Inventory

type Cursor struct {
	Start string
	End   string
}

func NewCursor(inventories Inventories) Cursor {
	if len(inventories) == 0 {
		return Cursor{}
	}
	return Cursor{
		Start: inventories[0].SellerSku,
		End:   inventories[len(inventories)-1].SellerSku,
	}
}

type InventoryRepository struct{}

func (r InventoryRepository) Save(ctx context.Context, db *bun.DB, inventories Inventories) error {
	_, err := db.NewInsert().
		Model(&inventories).
		On("CONFLICT (seller_sku) DO UPDATE").
		Set(strings.Join([]string{
			"asin = EXCLUDED.asin",
			"fnsku = EXCLUDED.fnsku",
			"condition = EXCLUDED.condition",
			"product_name = EXCLUDED.product_name",
			"quantity = EXCLUDED.quantity",
			"updated_at = current_timestamp",
		}, ",")).
		Exec(ctx)
	return err
}

func (r InventoryRepository) GetAll(ctx context.Context, db *bun.DB) (Inventories, error) {
	var inventories Inventories
	err := db.NewSelect().
		Model(&inventories).
		Order("seller_sku").
		Scan(ctx)
	return inventories, err
}

func (r InventoryRepository) GetBySellerSKU(ctx context.Context, db *bun.DB, skus []string) (Inventories, error) {
	var inventories Inventories
	if err := db.NewSelect().
		Model(&inventories).
		Where("seller_sku IN (?)", bun.In(skus)).
		Scan(ctx); err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r InventoryRepository) GetNextPage(ctx context.Context, db *bun.DB, cursor string, limit int) (Inventories, Cursor, error) {
	var inventories Inventories
	if err := db.NewSelect().
		Model(&inventories).
		Where("seller_sku > ?", cursor).
		Where("quantity > 0").
		Order("seller_sku ASC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, Cursor{}, err
	}
	return inventories, NewCursor(inventories), nil
}