package product

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type KeepaService struct{}

func (s KeepaService) UpdateRenderData(d amqp.Delivery) {

}
