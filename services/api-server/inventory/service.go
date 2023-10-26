package inventory

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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

func (i InventoryService) UpdatePricing() error {
	skus := []string{"4719512109847-N-1756-20231001"}
	query := url.Values{}
	query.Set("ids", strings.Join(skus, ","))
	query.Set("id_type", "Sku")
	URL, err := url.Parse(SpapiServiceURL)
	if err != nil {
		return err
	}
	URL.Path = "get-pricing"
	URL.RawQuery = query.Encode()
	res, err := http.Get(URL.String())
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
