package handler

import (
	"api-gateway/config"
	"api-gateway/models"
	"api-gateway/utils"
	"api-gateway/webResponse"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"messaging"
	"net/http"
	"user-service/api"
)

type UserHandler struct {
	Config          *config.RateLimitConfig
	RMQ             *messaging.RabbitMQConnection
	SendMessage     *api.SendingMessage
	ResponseHandler *webResponse.ResponseHandler
}

func NewUserHandler(cfg *config.RateLimitConfig, rmq *messaging.RabbitMQConnection, res *webResponse.ResponseHandler) *UserHandler {
	return &UserHandler{
		Config:          cfg,
		RMQ:             rmq,
		ResponseHandler: res,
		SendMessage:     api.NewSendingMessage(rmq),
	}
}

// Register handles user registration event-driven
func (h *UserHandler) Register(c echo.Context) error {
	// Rate Limit
	err := config.CheckRateLimit(c)
	if err != nil {
		return err
	}
	// bind & validate request
	var requestBody models.RegisterRequest
	if err := c.Bind(&requestBody); err != nil {
		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, "Invalid request format")
	}
	if err := requestBody.Validate(); err != nil {
		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, err.Error())
	}

	// Generate Correlation ID
	correlationID := utils.GenerateCorrelationID()
	logrus.Infof("Sending UserRegistered event | Correlation ID: %s | Payload: %+v", correlationID, requestBody)
	err = h.SendMessage.SendingToMessage("UserRegistered", correlationID, requestBody)
	if err != nil {
		return err
	}

	logrus.Infof("Waiting for UserRegisteredSuccess/UserRegisteredFailed response (Timeout: %v)", h.Config.RequestTimeout)
	return h.ResponseHandler.HandleEventResponse(
		c,
		false,
		http.StatusCreated,
		h.Config.RequestTimeout,
		"User registered successfully",
		"UserRegisteredSuccess",
		"UserRegisteredFailed",
	)

}

// Login handles user login event-driven
func (h *UserHandler) Login(c echo.Context) error {
	// Rate Limit
	err := config.CheckRateLimit(c)
	if err != nil {
		return err
	}
	// bind & validate request
	var requestBody models.LoginRequest
	if err := c.Bind(&requestBody); err != nil {
		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, "Invalid request format")
	}
	if err := requestBody.Validate(); err != nil {
		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, err.Error())
	}

	// Generate Correlation ID
	correlationID := utils.GenerateCorrelationID()

	err = h.SendMessage.SendingToMessage("UserLogin", correlationID, requestBody)
	if err != nil {
		return err
	}

	return h.ResponseHandler.HandleEventResponse(
		c,
		true,
		http.StatusCreated,
		h.Config.RequestTimeout,
		"User login successfully",
		"UserLoginSuccess",
		"UserLoginFailed",
	)
}
