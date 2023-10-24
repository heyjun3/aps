package inventory

import (
	"api-server/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateQuantity(t *testing.T) {
	db := test.CreateTestDBConnection()
	test.ResetModel(context.Background(), db, &Inventory{})
	inventories := []*Inventory{
		{SellerSku: "sku", TotalQuantity: 0, ProductName: "sku", Price: Ptr[int](100)},
		{SellerSku: "test", TotalQuantity: 0, ProductName: "test", Price: Ptr[int](200)},
	}
	if err := inventoryRepository.Save(context.Background(), db, inventories); err != nil {
		panic(err)
	}

	tests := []struct {
		name              string
		inventories       []*Inventory
		expectInventories []*Inventory
		wantErr           bool
	}{{
		name: "update quantity",
		inventories: []*Inventory{
			{SellerSku: "sku", TotalQuantity: 10},
			{SellerSku: "test", TotalQuantity: 20},
			{SellerSku: "test1", TotalQuantity: 20},
		},
		expectInventories: []*Inventory{
			{SellerSku: "sku", TotalQuantity: 10, ProductName: "sku", Price: Ptr[int](100)},
			{SellerSku: "test", TotalQuantity: 20, ProductName: "test", Price: Ptr[int](200)},
			{SellerSku: "test1", TotalQuantity: 20},
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InventoryService{}.UpdateQuantity(context.Background(), db, tt.inventories)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			inventories, err := inventoryRepository.GetAll(context.Background(), db)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, len(tt.expectInventories), len(inventories))
			for i, expect := range tt.expectInventories {
				assert.Equal(t, expect.Asin, inventories[i].Asin)
				assert.Equal(t, expect.SellerSku, inventories[i].SellerSku)
				assert.Equal(t, expect.FnSku, inventories[i].FnSku)
				assert.Equal(t, expect.Condition, inventories[i].Condition)
				assert.Equal(t, expect.ProductName, inventories[i].ProductName)
				assert.Equal(t, expect.LastUpdatedTime, inventories[i].LastUpdatedTime)
				assert.Equal(t, expect.TotalQuantity, inventories[i].TotalQuantity)
				assert.Equal(t, expect.Price, inventories[i].Price)
				assert.Equal(t, expect.Point, inventories[i].Point)
				assert.Equal(t, expect.LowestPrice, inventories[i].LowestPrice)
				assert.Equal(t, expect.LowestPoint, inventories[i].LowestPoint)
			}
		})
	}
}
