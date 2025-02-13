package messaging

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"os"
)

// PublishMessage message to RabbitMQ
func PublishMessage(queueName string, body string) error {
	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	if rabbitMQURI == "" {
		rabbitMQURI = "amqp://guest:guest@localhost:5672/"
	}

	// Connect to RabbitMQ
	conn, err := amqp091.Dial(rabbitMQURI)
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer func(conn *amqp091.Connection) {
		err := conn.Close()
		if err != nil {
			logrus.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
		return err
	}

	defer func(ch *amqp091.Channel) {
		err := ch.Close()
		if err != nil {
			logrus.Fatalf("Failed to close channel: %v", err)
		}
	}(ch)

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
		return err
	}

	err = ch.Publish("", q.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		logrus.Fatalf("Failed to publish a message: %v", err)
		return err
	}
	logrus.Println("Message published successfully")
	return nil
}
