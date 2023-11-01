package inventory

import (
	"api-server/spapi"
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Inventory struct {
	bun.BaseModel `bun:"table:inventories"`
	*spapi.Inventory
	Price       *int      `bun:"price"`
	Point       *int      `bun:"point"`
	LowestPrice *int      `bun:"lowest_price"`
	LowestPoint *int      `bun:"lowest_point"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewInventory(
	asin, fnSku, sellerSku, condition, lastUpdatedTime, productName string,
	totalQuantity, price, point, lowestPrice, lowestPoint int,
) *Inventory {
	return &Inventory{
		Inventory: &spapi.Inventory{
			Asin:            asin,
			FnSku:           fnSku,
			SellerSku:       sellerSku,
			Condition:       condition,
			LastUpdatedTime: lastUpdatedTime,
			ProductName:     productName,
			TotalQuantity:   totalQuantity,
		},
		Price:       &price,
		Point:       &point,
		LowestPrice: &lowestPrice,
		LowestPoint: &lowestPoint,
	}
}

func mergeTotalQuantity(base *Inventory, i *Inventory) *Inventory {
	base.TotalQuantity = i.TotalQuantity
	return base
}

func mergePriceAndPoints(base *Inventory, i *Inventory) *Inventory {
	base.Price = i.Price
	base.Point = i.Point
	return base
}

type Cursor struct {
	Start string
	End   string
}

func NewCursor(inventories []*Inventory) Cursor {
	if len(inventories) == 0 {
		return Cursor{}
	}
	return Cursor{
		Start: inventories[0].SellerSku,
		End:   inventories[len(inventories)-1].SellerSku,
	}
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
			"updated_at = current_timestamp",
		}, ",")).
		Exec(ctx)
	return err
}

func (r InventoryRepository) GetAll(ctx context.Context, db *bun.DB) ([]*Inventory, error) {
	var inventories []*Inventory
	err := db.NewSelect().
		Model(&inventories).
		Order("seller_sku").
		Scan(ctx)
	return inventories, err
}

func (r InventoryRepository) GetBySellerSKU(ctx context.Context, db *bun.DB, skus []string) ([]*Inventory, error) {
	var inventories []*Inventory
	if err := db.NewSelect().
		Model(&inventories).
		Where("seller_sku IN (?)", bun.In(skus)).
		Scan(ctx); err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r InventoryRepository) GetNextPage(ctx context.Context, db *bun.DB, cursor string, limit int) ([]*Inventory, Cursor, error) {
	var inventories []*Inventory
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
