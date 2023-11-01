package inventory

import (
	"api-server/spapi"
	"api-server/test"
	"context"
	"fmt"
	"strconv"
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
			{Inventory: &spapi.Inventory{Asin: "asin", SellerSku: "sku", Condition: "New"}, Price: Ptr[int](100)},
		},
		wantErr: false,
	}, {
		name: "save inventories",
		inventories: []*Inventory{
			{Inventory: &spapi.Inventory{Asin: "asin", SellerSku: "sku", Condition: "New"}},
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
		{Inventory: &spapi.Inventory{SellerSku: "sku", ProductName: "sku"}},
		{Inventory: &spapi.Inventory{SellerSku: "test", ProductName: "test"}},
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

func TestGetNextPage(t *testing.T) {
	db := test.CreateTestDBConnection()
	test.ResetModel(context.Background(), db, &Inventory{})
	repo := InventoryRepository{}
	seed := make([]*Inventory, 100)
	for i := range seed {
		seed[i] = &Inventory{Inventory: &spapi.Inventory{SellerSku: strconv.Itoa(i + 1), TotalQuantity: 10}}
	}
	if err := repo.Save(context.Background(), db, seed); err != nil {
		panic(err)
	}

	type args struct {
		cursor string
		limit  int
	}
	type want struct {
		cursor Cursor
		count  int
		first  Inventory
		last   Inventory
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "get next page",
			args:    args{cursor: "10", limit: 20},
			want:    want{cursor: Cursor{Start: "100", End: "28"}, first: Inventory{Inventory: &spapi.Inventory{SellerSku: "100"}}, last: Inventory{Inventory: &spapi.Inventory{SellerSku: "28"}}, count: 20},
			wantErr: false,
		},
		{
			name:    "no page",
			args:    args{cursor: "998", limit: 100},
			want:    want{count: 0, cursor: Cursor{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventories, cursor, err := repo.GetNextPage(context.Background(), db, tt.args.cursor, tt.args.limit)

			fmt.Println(inventories, cursor, err)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.count, len(inventories))
				if tt.want.count > 0 {
					assert.Equal(t, tt.want.first.SellerSku, inventories[0].SellerSku)
					assert.Equal(t, tt.want.last.SellerSku, inventories[len(inventories)-1].SellerSku)
				}
				assert.Equal(t, tt.want.cursor, cursor)
			}
		})
	}
}
