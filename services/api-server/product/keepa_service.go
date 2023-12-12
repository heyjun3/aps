package product

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type KeepaService struct{}

func (s KeepaService) Consume() {
	conn, err := amqp.Dial(os.Getenv("MQ_DSN"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("chart", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare a queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Recived a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
