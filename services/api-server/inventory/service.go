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
	mergeInventories := MergeInventories(inventoriesInDB, inventories, mergeTotalQuantity)
	return inventoryRepository.Save(ctx, db, mergeInventories)
}

func MergeInventories(dst []*Inventory, src []*Inventory, fn func(*Inventory, *Inventory) *Inventory) []*Inventory {
	inventoryMap := make(map[string]*Inventory)
	for _, inventory := range dst {
		inventoryMap[inventory.SellerSku] = inventory
	}
	mergedInventories := make([]*Inventory, 0, len(src))
	for _, inventory := range src {
		inventoryForMap := inventoryMap[inventory.SellerSku]
		if inventoryForMap == nil {
			mergedInventories = append(mergedInventories, inventory)
			continue
		}
		mergedInventories = append(mergedInventories, fn(inventoryForMap, inventory))
	}
	return mergedInventories
}
