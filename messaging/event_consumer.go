package messaging

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// ResponseChannel is a channel to send response to API Gateway
var ResponseChannel = make(chan string)

// ConsumeMessage message from RabbitMQ
func ConsumeMessage(queueName string, handler func([]byte, chan string)) {
	rabbitMQURI := os.Getenv("RABBITMQ_URI")

	conn, err := amqp091.Dial(rabbitMQURI)
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("Failed to consume messages: %v", err)
	}

	logrus.Printf("Listening for messages from queue: %s", queueName)

	go func() {
		for d := range msgs {
			logrus.Infof("Received message from %s: %s", queueName, d.Body)
			handler(d.Body, ResponseChannel)
		}
	}()
}
