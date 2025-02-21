package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"time"
	"user-service/core/models"
)

// ConsumeEvent listens for specific event types
func ConsumeEvent(rmq *RabbitMQConnection, serviceName string, eventNames []string, handler func(event models.Event)) {
	ch, err := rmq.GetChannel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
		return
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_queue", serviceName),
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

	for _, eventName := range eventNames {
		err = ch.QueueBind(q.Name, eventName, "events_exchange", false, nil)
		if err != nil {
			logrus.Fatalf("Failed to bind queue %s to event %s: %v", q.Name, eventName, err)
			return
		}
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("Failed to consume messages: %v", err)
		return
	}

	go func() {
		for d := range msgs {
			var event models.Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				logrus.Errorf("Failed to parse event data: %v", err)
				continue
			}
			logrus.Infof("[RabbitMQ] Event received: %s | CorrelationID: %s", event.EventType, event.CorrelationID)

			found := false
			for _, expectedEvent := range eventNames {
				if event.EventType == expectedEvent {
					found = true
					break
				}
			}

			if found {
				handler(event)
				err = d.Ack(false)
				if err != nil {
					logrus.Errorf("Failed to acknowledge message: %v", err)
				}
			} else {
				logrus.Warnf("Received unexpected event: %s", event.EventType)
			}
		}
	}()
}

// WaitForEvent waits for the specified events to arrive within a given timeout
func WaitForEvent(rmq *RabbitMQConnection, timeout time.Duration, serviceName string, eventNames ...string) (models.Event, error) {
	ch, err := rmq.GetChannel()
	if err != nil {
		logrus.Errorf("Failed to open RabbitMQ channel: %v", err)
		return models.Event{}, err
	}
	defer ch.Close()

	msgChannel := make(chan models.Event, 1)
	errChan := make(chan error, 1)

	queue, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("Failed to declare queue: %v", err)
		return models.Event{}, err
	}

	for _, eventName := range eventNames {
		err = ch.QueueBind(queue.Name, eventName, "events_exchange", false, nil)
		if err != nil {
			logrus.Errorf("Failed to bind queue to event %s: %v", eventName, err)
			return models.Event{}, err
		}
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Errorf("Failed to consume messages: %v", err)
		return models.Event{}, err
	}

	go func() {
		for msg := range msgs {
			var event models.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				logrus.Errorf("Failed to parse event data: %v", err)
				errChan <- err
				continue
			}

			logrus.Infof("[RabbitMQ] Event received: %s | CorrelationID: %s", event.EventType, event.CorrelationID)

			for _, expectedEvent := range eventNames {
				if event.EventType == expectedEvent {
					msg.Ack(false)
					msgChannel <- event
					return
				}
			}
		}
	}()

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
func (rmq *RabbitMQConnection) GetChannel() (*amqp091.Channel, error) {
	if rmq.channel == nil || rmq.channel.IsClosed() {
		ch, err := rmq.conn.Channel()
		if err != nil {
			return nil, err
		}
		rmq.channel = ch
	}
	return rmq.channel, nil
}
