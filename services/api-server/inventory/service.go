package inventory

import (
	"context"

	"github.com/uptrace/bun"
)

type InventoryService struct{}

func (i InventoryService) UpdateQuantity(ctx context.Context, db *bun.DB, inventories []*Inventory) error {
	sellerSkus := []string{}
	for _, inventory := range inventories {
		sellerSkus = append(sellerSkus, inventory.SellerSku)
	}

	inventoriesInDB, err := inventoryRepository.GetBySellerSKU(ctx, db, sellerSkus)
	if err != nil {
		return err
	}
	inventoryMap := make(map[string]*Inventory)
	for _, inventory := range inventoriesInDB {
		inventoryMap[inventory.SellerSku] = inventory
	}

	updateInventories := []*Inventory{}
	for _, inventory := range inventories {
		inventoryInDB := inventoryMap[inventory.SellerSku]
		if inventoryInDB == nil {
			updateInventories = append(updateInventories, inventory)
			continue
		}
		inventoryInDB.SetTotalQuantity(inventory.TotalQuantity)
		updateInventories = append(updateInventories, inventoryInDB)
	}
	return inventoryRepository.Save(ctx, db, updateInventories)
}
