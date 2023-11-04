package inventory

import (
	"context"

	"github.com/uptrace/bun"
)

type InventoryService struct{}

func (i InventoryService) UpdateQuantity(ctx context.Context, db *bun.DB, inventories Inventories) error {
	sellerSkus := inventories.Skus()
	inventoriesInDB, err := inventoryRepository.GetBySellerSKU(ctx, db, sellerSkus)
	if err != nil {
		return err
	}
	mergeInventories := MergeInventories(inventoriesInDB, inventories, mergeTotalQuantity)
	return inventoryRepository.Save(ctx, db, mergeInventories)
}

func MergeInventories(dst Inventories, src Inventories, fn func(*Inventory, *Inventory) *Inventory) Inventories {
	inventoryMap := dst.Map()
	mergedInventories := make(Inventories, 0, len(src))
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
