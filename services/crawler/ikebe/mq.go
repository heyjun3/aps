package ikebe

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func createMQConnection(name string) {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ")
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open channel")
		return
	}

	q, err := ch.QueueDeclare("mws", true, false, false, false, nil)
	if err != nil {
		log.Fatalln("Failed to declare a queue")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	body := "hello"
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
}