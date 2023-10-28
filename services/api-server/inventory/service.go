package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"
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

type GetPricingResponse struct {
	Paylaod []GetPricingPayload `json:"payload"`
}

type GetPricingPayload struct {
	Status    string            `json:"status"`
	SellerSKU string            `json:"SellerSKU"`
	Product   GetPricingProduct `json:"Product"`
}

type GetPricingProduct struct {
	Offers []Offers `json:"Offers"`
}

type Offers struct {
	BuyingPrice BuyingPrice `json:"BuyingPrice"`
}

type BuyingPrice struct {
	ListingPrice Price  `json:"ListingPrice"`
	Points       Points `json:"Points"`
}

type Price struct {
	CurrencyCode string  `json:"CurrencyCode"`
	Amount       float64 `json:"Amount"`
}

type Points struct {
	PointsNumber int64 `json:"PointsNumber"`
}

func (i InventoryService) UpdatePricing() error {
	skus := []string{"4719512109847-N-1756-20231001", "4515260019373-B-2500-20231025"}
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
	var getPricingResponse GetPricingResponse
	if err := json.Unmarshal(body, &getPricingResponse); err != nil {
		slog.Error("err", err)
		return err
	}
	fmt.Println(getPricingResponse.Paylaod[0].Product.Offers[0].BuyingPrice)
	return nil
}
