package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
	"user-service/core/models"
)

// ConsumeEvent listens for specific event types
func ConsumeEvent(rmq *RabbitMQConnection, eventName string, handler func(event models.Event)) {
	ch, err := rmq.conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
		return
	}

	// Declare queue for event type (event name)
	q, err := ch.QueueDeclare(
		eventName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
		return
	}

	// Bind queue ke exchange
	err = ch.QueueBind(q.Name, eventName, "events_exchange", false, nil)
	if err != nil {
		logrus.Fatalf("Failed to bind queue: %v", err)
		return
	}

	// consume messages from queue
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("Failed to consume messages: %v", err)
		return
	}

	logrus.Infof("[RabbitMQ] Listening for event: %s on queue: %s", eventName, q.Name)

	go func() {
		for d := range msgs {
			var event models.Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				logrus.Errorf("Failed to parse event data: %v", err)
				continue
			}

			handler(event)

			err = d.Ack(false)
			if err != nil {
				logrus.Errorf("Failed to acknowledge message: %v", err)
			}
		}
	}()
}

// WaitForEvent waits for the specified events to arrive within a given timeout
func WaitForEvent(rmq *RabbitMQConnection, timeout time.Duration, eventNames ...string) (models.Event, error) {
	ch, err := rmq.conn.Channel()
	if err != nil {
		logrus.Errorf("Failed to open RabbitMQ channel: %v", err)
		return models.Event{}, err
	}
	defer ch.Close()

	msgChannel := make(chan models.Event, 1)
	errChan := make(chan error, 1)

	// consume messages from queue
	for _, eventName := range eventNames {
		go func(eventName string) {
			msgs, err := ch.Consume(eventName, "", false, false, false, false, nil)
			if err != nil {
				errChan <- fmt.Errorf("error subscribing to event %s: %v", eventName, err)
				return
			}
			for msg := range msgs {
				logrus.Infof("Received message from queue: %s | Body: %s", eventName, string(msg.Body))

				var event models.Event
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					errChan <- fmt.Errorf("error unmarshalling event: %v", err)
					return
				}

				logrus.Infof("Event received: %s | CorrelationID: %s", event.EventType, event.CorrelationID)
				msgChannel <- event
				return
			}
		}(eventName)
	}

	select {
	case event := <-msgChannel:
		return event, nil

	case err := <-errChan:
		return models.Event{}, err

	case <-time.After(timeout):
		logrus.Errorf("Event timeout while waiting for: %v", eventNames)
		return models.Event{}, fmt.Errorf("timeout while waiting for event: %v", eventNames)
	}
}
