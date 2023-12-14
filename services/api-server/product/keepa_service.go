package product

import (
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"api-server/spapi/price/competitive"
)

type KeepaService struct{}

func (s KeepaService) UpdateRenderData(d amqp.Delivery) {
	var res competitive.GetCompetitivePricingResponse
	if err := json.Unmarshal(d.Body, &res); err != nil {
		slog.Error("json unmarshal error", err)
		return
	}
	// prices := res.LandedPrices()
}
