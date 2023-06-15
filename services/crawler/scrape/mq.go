package scrape

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"crawler/config"
)

type RabbitMQClient interface {
	Publish([]byte) error
}

type MQClient struct {
	dsn       string
	queueName string
}

func NewMQClient(dsn, name string) MQClient {
	return MQClient{dsn: dsn, queueName: name}
}

func (mq MQClient) createMQConnection() (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(mq.dsn)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open channel", err)
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(mq.queueName, true, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to declare a queue", err)
		return nil, nil, err
	}

	return ch, conn, err
}

func (mq MQClient) CreateConsumer() (<-chan amqp.Delivery, *amqp.Channel, error) {
	ch, _, err := mq.createMQConnection()

	if err != nil {
		ch.Close()
		return nil, nil, err
	}
	msgs, err := ch.Consume(mq.queueName, "", false, false, false, false, nil)
	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	return msgs, ch, nil
}

func (mq MQClient) Publish(message []byte) error {
	ch, conn, err := mq.createMQConnection()
	defer conn.Close()
	defer ch.Close()

	if err != nil {
		logger.Error("create connection error", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = ch.PublishWithContext(ctx, "", mq.queueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
	return err
}

func MoveMessages(srcQueue, dstQueue string) {
	srcClient := NewMQClient(config.MQDsn, srcQueue)
	msgs, ch, err := srcClient.CreateConsumer()
	defer ch.Close()
	if err != nil {
		logger.Error("error", err)
	}
	dstClient := NewMQClient(config.DstMQDsn, dstQueue)
	for d := range msgs {
		logger.Info(string(d.Body))
		err := dstClient.Publish(d.Body)
		if err != nil {
			logger.Error("publish error", err)
			d.Nack(true, true)
			return
		}
		d.Ack(true)
	}
}
