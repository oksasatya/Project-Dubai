package webResponse

import (
	"api-gateway/config"
	"api-gateway/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"messaging"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type ResponseHandler struct {
	RMQ *messaging.RabbitMQConnection
}

// NewResponseHandler creates a new instance of ResponseHandler
func NewResponseHandler(rmq *messaging.RabbitMQConnection) *ResponseHandler {
	if rmq == nil {
		logrus.Fatal("NewResponseHandler: RabbitMQ connection is nil!")
	}
	return &ResponseHandler{
		RMQ: rmq,
	}
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// ResponseJson is a utility function that sends a JSON response to the client
func ResponseJson(c echo.Context, status int, payload interface{}, message string) error {
	response := Response{
		Meta: Meta{
			Message: message,
			Code:    status,
			Status:  getStatusText(status),
		},
		Data: payload,
	}
	return c.JSON(status, response)
}

func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	} else if status >= 400 && status < 500 {
		return "fail"
	} else {
		return "error"
	}
}

// HandleEventResponse handles event-based response waiting and processing
func (h *ResponseHandler) HandleEventResponse(c echo.Context, generateToken bool, statusCode int, timeout time.Duration, message string, eventName ...string) error {
	if h == nil {
		logrus.Fatal("HandleEventResponse: ResponseHandler is nil!")
		return ResponseJson(c, http.StatusInternalServerError, nil, "Internal Server Error: ResponseHandler is nil")
	}

	if h.RMQ == nil {
		logrus.Fatal("HandleEventResponse: RabbitMQ connection is nil!")
		return ResponseJson(c, http.StatusInternalServerError, nil, "Internal Server Error: RabbitMQ connection is nil")
	}

	responseEvent, err := messaging.WaitForEvent(h.RMQ, timeout, "api-gateway", eventName...)
	if err != nil {
		logrus.Errorf("Event timeout while waiting for: %s", eventName)
		return ResponseJson(c, http.StatusGatewayTimeout, nil, "Request timed out waiting for response")
	}

	logrus.Infof("Received event: %s | CorrelationID: %s", eventName, responseEvent.CorrelationID)
	ctx := c.Request().Context()

	var jsonResponse map[string]interface{}
	if payloadStr, ok := responseEvent.Payload.(string); ok {
		//logrus.Warnf("Payload is a string: %s", payloadStr)
		if json.Valid([]byte(payloadStr)) {
			if err := json.Unmarshal([]byte(payloadStr), &jsonResponse); err != nil {
				logrus.Errorf("Failed to parse event payload: %v", err)
				return ResponseJson(c, http.StatusInternalServerError, nil, "Failed to parse response")
			}
		} else {
			return ResponseJson(c, http.StatusConflict, nil, payloadStr)
		}
	} else if payloadMap, ok := responseEvent.Payload.(map[string]interface{}); ok {
		jsonResponse = payloadMap
	} else {
		logrus.Errorf("Unexpected payload type: %T", responseEvent.Payload)
		return ResponseJson(c, http.StatusInternalServerError, nil, "Unexpected event payload format")
	}

	for _, expectedEvent := range eventName {
		if responseEvent.EventType == expectedEvent {
			//logrus.Infof("Received expected event: %s", expectedEvent)
			delete(jsonResponse, "password")
			if generateToken {
				userID, _ := jsonResponse["id"].(string)
				userEmail, _ := jsonResponse["email"].(string)
				userRole, _ := jsonResponse["role"].(string)
				token, err := utils.GenerateToken(userID, userEmail, userRole)
				if err != nil {
					logrus.Error("Failed to generate token")
					return ResponseJson(c, http.StatusInternalServerError, nil, "Failed to generate token")
				}

				// Simpan token ke Redis
				ttlHoursStr := os.Getenv("JWT_EXPIRATION_TIME")
				ttlHours := 72
				if ttlHoursStr != "" {
					var err error
					ttlHours, err = strconv.Atoi(ttlHoursStr)
					if err != nil {
						logrus.Errorf("Invalid JWT_EXPIRATION_TIME value, using default: %v", err)
					}
				}

				// Redis client
				storeToken := config.StoreTokenInRedis(ctx, token, userID, userRole, ttlHours)
				if storeToken != nil {
					logrus.Errorf("Failed to store token in Redis: %v", storeToken)
					return ResponseJson(c, http.StatusInternalServerError, nil, "Failed to store token in Redis")
				}

				logrus.Infof("Token stored in Redis: %s", token)

				jsonResponse["token"] = token
			}

			return ResponseJson(c, statusCode, jsonResponse, message)
		}
	}

	logrus.Errorf("Unexpected event received: %s", responseEvent.EventType)
	return ResponseJson(c, http.StatusInternalServerError, nil, "Unexpected event received")
}
