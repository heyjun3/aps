package product

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/bun"

	"api-server/spapi/price/competitive"
)

type KeepaService struct {
	repository KeepaRepository
}

func NewKeepaService(db *bun.DB) *KeepaService {
	return &KeepaService{
		repository: KeepaRepository{
			DB: db,
		},
	}
}

func (s KeepaService) UpdateRenderData(d amqp.Delivery) {
	var res competitive.GetCompetitivePricingResponse
	if err := json.Unmarshal(d.Body, &res); err != nil {
		slog.Error("json unmarshal error", err)
		return
	}
	products := res.LandedPrices()
	keepas, err := s.repository.GetByAsins(context.Background(), products.Asins())
	if err != nil {
		slog.Error("failed get keepa", err)
		return
	}

}
