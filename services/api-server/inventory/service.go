package inventory

import (
	"context"
	"log/slog"

	"api-server/spapi"
	"api-server/spapi/price"
)

type InventoryService struct {
	spapiClient            *spapi.SpapiClient
	inventoryRepository    InventoryRepository
	currentPriceRepository PriceRepository[*CurrentPrice]
	lowestPriceRepository  PriceRepository[*LowestPrice]
	desiredPriceRepository PriceRepository[*DesiredPrice]
}

func NewInventoryService(spapiServiceURL string) (*InventoryService, error) {
	spapiClient, err := spapi.NewSpapiClient(spapiServiceURL)
	if err != nil {
		return nil, err
	}
	return &InventoryService{
		spapiClient:            spapiClient,
		inventoryRepository:    InventoryRepository{},
		currentPriceRepository: PriceRepository[*CurrentPrice]{},
		lowestPriceRepository:  PriceRepository[*LowestPrice]{},
		desiredPriceRepository: PriceRepository[*DesiredPrice]{},
	}, nil
}

func (s InventoryService) refreshInventories() error {
	ctx := context.Background()
	var nextToken string
	for {
		res, err := s.spapiClient.InventorySummaries(nextToken)
		if err != nil {
			return err
		}
		inventories := Inventories{}
		for _, inventory := range res.Payload.InventorySummaries {
			iv, err := NewInventoryFromInventory(inventory)
			if err != nil {
				slog.Warn(err.Error(), "struct", *inventory, "sku", inventory.SellerSku)
				continue
			}
			inventories = append(inventories, iv)
		}
		if err := s.inventoryRepository.Save(ctx, db, inventories); err != nil {
			slog.Error("failed save inventories", err)
			return err
		}
		slog.Info("logging next token", "nextToken", res.Pagination.NextToken)
		nextToken = res.Pagination.NextToken
		if nextToken == "" {
			slog.Info("break loop")
			break
		}
	}
	return nil
}

func (s InventoryService) refreshPricing() error {
	var inventories Inventories
	var cursor Cursor
	var err error
	for {
		inventories, cursor, err = s.inventoryRepository.GetNextPage(context.Background(), db, cursor.End, 20)
		if err != nil {
			slog.Error("error", "detail", err)
			return err
		}
		if len(inventories) == 0 || inventories == nil {
			return nil
		}

		res, err := s.spapiClient.GetPricing(inventories.Skus())
		if err != nil {
			slog.Error("error", "detail", err)
			return err
		}
		currentPrices := make(CurrentPrices, 0, len(inventories))
		lowestPrices := make(LowestPrices, 0, len(inventories))
		for _, response := range res.Responses {
			offers := response.Body.Payload.Offers
			if len(offers) == 0 {
				continue
			}
			sku := response.Body.Payload.SKU
			currents := offers.FilterCondition(price.Condition{MyOffer: Ptr(true)})
			if len(currents) > 0 {
				amount := currents[0].Price.Amount
				point := currents[0].Points.PointsNumber
				price, err := NewCurrentPrice(*sku, amount, point)
				if err == nil {
					currentPrices = append(currentPrices, price)
				}
			}

			lowest := getLowestPricingFromOffers(offers)
			if lowest != nil {
				amount := lowest.Price.Amount
				point := lowest.Points.PointsNumber
				price, err := NewLowestPrice(*sku, amount, point)
				if err != nil {
					slog.Error(err.Error(), "sku", sku)
					continue
				}
				lowestPrices = append(lowestPrices, price)
			}
		}
		if err := s.lowestPriceRepository.Save(context.Background(), db, lowestPrices); err != nil {
			slog.Error("error", "detail", err)
			return err
		}
		if err := s.currentPriceRepository.Save(context.Background(), db, currentPrices); err != nil {
			slog.Error("error", "detail", err)
			return err
		}
	}
}

func getLowestPricingFromOffers(offers price.Offers) *price.Offer {
	buyBoxWinnerAndFullfilledByAmazon := offers.FilterCondition(price.Condition{
		IsFullfilledByAmazon: Ptr(true),
		IsBuyBoxWinner:       Ptr(true),
	})
	if len(buyBoxWinnerAndFullfilledByAmazon) > 0 {
		return &buyBoxWinnerAndFullfilledByAmazon[0]
	}
	buyBox := offers.FilterCondition(price.Condition{
		IsBuyBoxWinner: Ptr(true),
	})
	if len(buyBox) > 0 {
		return &buyBox[0]
	}
	return offers.Lowest()
}

func (s InventoryService) getInventories() (Inventories, error) {
	ctx := context.Background()
	condition := Condition{Quantity: Ptr(0), IsNotOnlyLowestPrice: true}
	inventories, err := s.inventoryRepository.GetByCondition(ctx, db, condition)
	if err != nil {
		slog.Error("get inventories error", "detail", err)
		return nil, err
	}
	return inventories, nil
}

type updatePricingDTO struct {
	Sku          string  `json:"sku"`
	Price        float64 `json:"price"`
	PercnetPoint float64 `json:"percentPoint"`
}
type updatePricingDTOS []updatePricingDTO

func (d updatePricingDTOS) skus() []string {
	skus := make([]string, 0, len(d))
	for _, dto := range d {
		skus = append(skus, dto.Sku)
	}
	return skus
}

func (s InventoryService) updatePricing(dtos updatePricingDTOS) error {
	skus := dtos.skus()
	condition := Condition{Skus: skus}
	inventories, err := s.inventoryRepository.GetByCondition(context.Background(), db, condition)
	if err != nil {
		return err
	}
	m := inventories.Map()

	prices := make(DesiredPrices, 0, len(dtos))
	for _, dto := range dtos {
		d := dto
		inventory := m[d.Sku]
		if inventory == nil {
			continue
		}
		lowestPrice := inventory.LowestPrice
		if lowestPrice == nil {
			continue
		}

		price, err := NewDesiredPrice(&d.Sku, &d.Price, &d.PercnetPoint, *lowestPrice)
		if err != nil {
			return err
		}
		prices = append(prices, price)
	}
	if err := s.desiredPriceRepository.Save(context.Background(), db, prices); err != nil {
		return err
	}

	for _, price := range prices {
		p := price
		if err := s.spapiClient.UpdatePricing(p); err != nil {
			return err
		}
	}

	if err := s.spapiClient.UpdatePoints(prices); err != nil {
		return err
	}
	return nil
}
