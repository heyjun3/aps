package inventory

import (
	"api-server/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Ptr[T any](t T) *T {
	return &t
}

func TestSaveInventories(t *testing.T) {
	tests := []struct {
		name        string
		inventories []*Inventory
		wantErr     bool
	}{{
		name: "save inventories",
		inventories: []*Inventory{
			{Asin: "asin", SellerSku: "sku", Condition: "New", Price: Ptr[int](100)},
		},
		wantErr: false,
	}, {
		name: "save inventories",
		inventories: []*Inventory{
			{Asin: "asin", SellerSku: "sku", Condition: "New"},
		},
		wantErr: false,
	}}

	db := test.CreateTestDBConnection()
	test.ResetModel(context.Background(), db, &Inventory{})
	repo := InventoryRepository{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(context.Background(), db, tt.inventories)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBySellerSKU(t *testing.T) {
	db := test.CreateTestDBConnection()
	test.ResetModel(context.Background(), db, &Inventory{})
	repo := InventoryRepository{}
	inventories := []*Inventory{
		{SellerSku: "sku", ProductName: "sku"},
		{SellerSku: "test", ProductName: "test"},
	}
	if err := repo.Save(context.Background(), db, inventories); err != nil {
		panic(err)
	}

	tests := []struct {
		name        string
		skus        []string
		inventories []*Inventory
		wantErr     bool
	}{{
		name:        "get by seller sku",
		skus:        []string{"test", "sku"},
		inventories: inventories,
		wantErr:     false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventories, err := repo.GetBySellerSKU(context.Background(), db, tt.skus)

			assert.Equal(t, tt.inventories, inventories)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
