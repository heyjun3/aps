package inventory

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/uptrace/bun"

	inventory "api-server/spapi/inventory"
)

func Ptr[T any](v T) *T {
	return &v
}

func ValidateNilFieldsOfStruct[T any](value *T) (*T, error) {
	v := reflect.ValueOf(*value)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == nil {
			return nil, errors.New("nil field in Struct")
		}
	}
	return value, nil
}

type Inventory struct {
	bun.BaseModel       `bun:"table:inventories"`
	Asin                *string       `bun:"asin"`
	FnSku               *string       `bun:"fnsku"`
	SellerSku           *string       `bun:"seller_sku,pk"`
	Condition           *string       `bun:"condition"`
	LastUpdatedTime     *string       `bun:"-"`
	ProductName         *string       `bun:"product_name"`
	TotalQuantity       *int          `bun:"quantity"`
	FulfillableQuantity *int          `bun:"fulfillable_quantity"`
	CurrentPrice        *CurrentPrice `bun:"rel:has-one,join:seller_sku=seller_sku"`
	LowestPrice         *LowestPrice  `bun:"rel:has-one,join:seller_sku=seller_sku"`
	DesiredPrice        *DesiredPrice `bun:"rel:has-one,join:seller_sku=seller_sku"`
	CreatedAt           time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt           time.Time     `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewInventory(
	asin, fnSku, sellerSku, condition, productName string,
	totalQuantity, fulfillableQuantity int,
) *Inventory {
	return &Inventory{
		Asin:                &asin,
		FnSku:               &fnSku,
		SellerSku:           &sellerSku,
		Condition:           &condition,
		ProductName:         &productName,
		TotalQuantity:       &totalQuantity,
		FulfillableQuantity: &fulfillableQuantity,
	}
}

func NewInventoryFromInventory(iv *inventory.Inventory) (*Inventory, error) {
	value, err := ValidateNilFieldsOfStruct[inventory.Inventory](iv)
	if err != nil {
		return nil, err
	}
	inventory := NewInventory(
		*value.Asin,
		*value.FnSku,
		*value.SellerSku,
		*value.Condition,
		*value.ProductName,
		*value.TotalQuantity,
		*value.InventoryDetails.FulfillableQuantity,
	)
	return inventory, nil
}

type Inventories []*Inventory

func (i *Inventories) Skus() []string {
	skus := make([]string, 0, len(*i))
	for _, iv := range *i {
		skus = append(skus, *iv.SellerSku)
	}
	return skus
}

func (i Inventories) Map() map[string]*Inventory {
	m := make(map[string]*Inventory)
	for _, iv := range i {
		m[*iv.SellerSku] = iv
	}
	return m
}

type Cursor struct {
	Start string
	End   string
}

func NewCursor(inventories Inventories) Cursor {
	if len(inventories) == 0 {
		return Cursor{}
	}
	return Cursor{
		Start: *inventories[0].SellerSku,
		End:   *inventories[len(inventories)-1].SellerSku,
	}
}

type Condition struct {
	Quantity             *int
	IsNotOnlyLowestPrice bool
	Skus                 []string
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

func (r InventoryRepository) GetByCondition(ctx context.Context, db *bun.DB, condition Condition) (Inventories, error) {
	var inventories Inventories
	query := db.NewSelect().
		Model(&inventories).
		Relation("CurrentPrice").
		Relation("LowestPrice").
		Relation("DesiredPrice").
		Order("seller_sku")

	if len(condition.Skus) > 0 {
		query.Where("inventory.seller_sku IN (?)", bun.In(condition.Skus))
	}

	if condition.Quantity != nil {
		query.Where("quantity > ?", *condition.Quantity)
	}
	if condition.IsNotOnlyLowestPrice {
		query.WhereGroup("AND", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("current_price.amount != lowest_price.amount").
				WhereOr("current_price.point != lowest_price.point")
		})
	}

	err := query.Scan(ctx)
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
