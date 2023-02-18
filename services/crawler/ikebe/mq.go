package ikebe

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	publish([]byte) error
}

type MQClient struct {
	dsn string
	queueName string
}

func NewMQClient(dsn, name string) MQClient{
	return MQClient{dsn: dsn, queueName: name}
}

func (mq MQClient) createMQConnection() (*amqp.Channel, error){
	conn, err := amqp.Dial(mq.dsn)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open channel", err)
		return nil, err
	}

	_, err = ch.QueueDeclare(mq.queueName, true, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to declare a queue", err)
		return nil, err
	}

	return ch, err
}

func (mq MQClient) publish(message []byte) error {
	ch, err := mq.createMQConnection()
	if err != nil {
		logger.Error("create connection error", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	err = ch.PublishWithContext(ctx, "", mq.queueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
	return err
}

func (mq MQClient) batchPublish(messages ...[]byte) error {
	ch, err := mq.createMQConnection()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, message := range messages {
		err = ch.PublishWithContext(ctx, "", mq.queueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
		if err != nil {
			return err
		}
	}
	return err
}

type MWSSchema struct {
	Filename string `json:"filename"`
	Jan string `json:"jan"`
	Price int64 `json:"price"`
	URL string `json:"url"`
}

func NewMWSSchema(filename, jan, url string, price int64) *MWSSchema{
	return &MWSSchema{
		Filename: filename,
		Jan: jan,
		URL: url,
		Price: price,
	}
}