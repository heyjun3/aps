package product

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/bun"

	"api-server/spapi/price/competitive"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

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
	defer d.Ack(true)

	var res competitive.GetCompetitivePricingResponse
	if err := json.Unmarshal(d.Body, &res); err != nil {
		logger.Error("json unmarshal error", err)
		return
	}
	renderData := ConvertLandedProducts(res.LandedPrices())

	keepas, err := s.repository.GetByAsins(context.Background(), renderData.Asins())
	if err != nil {
		logger.Error("failed get keepa", "err", err)
		return
	}
	if err := s.repository.Save(context.Background(), keepas.UpdateRenderData(renderData)); err != nil {
		logger.Error("failed save keepa", "err", err)
		return
	}
	logger.Info("update render data is done")
}
