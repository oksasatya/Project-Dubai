package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"messaging"
	"time"
	"user-service/core/models"
)

type SendingMessage struct {
	Rmq *messaging.RabbitMQConnection
}

func NewSendingMessage(rmq *messaging.RabbitMQConnection) *SendingMessage {
	return &SendingMessage{
		Rmq: rmq,
	}
}

// SendingToMessage is a function to send message to message broker
func (s *SendingMessage) SendingToMessage(eventType string, correlationID string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	// Create Event Object
	event := models.Event{
		EventType:     eventType,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
		Payload:       json.RawMessage(payloadBytes),
	}

	// Serialize Event
	eventJSON, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Failed to marshal event: %v", err)
	}

	// Publish Event
	err = s.Rmq.PublishEvent(eventType, eventJSON)
	if err != nil {
		logrus.Errorf("Failed to publish event: %v", err)
	}
	return err
}
