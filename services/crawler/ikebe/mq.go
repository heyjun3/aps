package ikebe

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	createMQConnection() (*amqp.Channel, error)
	publish([]byte) error
	batchPublish([][]byte) error
}

type MQClient struct {
	dsn string
	queueName string
}

func NewMQClient(dsn, name string) *MQClient{
	return &MQClient{dsn: dsn, queueName: name}
}

func (mq *MQClient)createMQConnection() (*amqp.Channel, error){
	conn, err := amqp.Dial(mq.dsn)
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open channel")
		return nil, err
	}

	_, err = ch.QueueDeclare(mq.queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("Failed to declare a queue")
		return nil, err
	}

	return ch, err
}

func (mq *MQClient) publish(message []byte) error {
	ch, err := mq.createMQConnection()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	err = ch.PublishWithContext(ctx, "", mq.queueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
	return err
}

func (mq *MQClient) batchPublish(messages ...[]byte) error {
	ch, err := mq.createMQConnection()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	ctx := context.Background()
	for _, message := range messages {
		err = ch.PublishWithContext(ctx, "", mq.queueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
		if err != nil {
			log.Fatalln(err)
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