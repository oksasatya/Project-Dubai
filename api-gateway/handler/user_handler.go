package handler

import (
	"api-gateway/config"
	"api-gateway/models"
	"api-gateway/utils"
	"api-gateway/webResponse"
	"errors"
	"github.com/go-playground/validator/v10"
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
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formatterErrors := utils.FormatValidationError(&requestBody, validationErrors)
			return webResponse.ResponseJson(c, http.StatusBadRequest, nil, formatterErrors)
		}
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
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formatterErrors := utils.FormatValidationError(&requestBody, validationErrors)
			return webResponse.ResponseJson(c, http.StatusBadRequest, nil, formatterErrors)
		}
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
		http.StatusAccepted,
		h.Config.RequestTimeout,
		"User login successfully",
		"UserLoginSuccess",
		"UserLoginFailed",
	)
}

// GetProfile handles user testing
func (h *UserHandler) GetProfile(c echo.Context) error {
	err := config.CheckRateLimit(c)
	if err != nil {
		return err
	}

	claims, ok := c.Get("user").(*utils.JWTCustomClaims)
	if !ok || claims == nil {
		logrus.Error("Invalid claims type or nil claims")
		return webResponse.ResponseJson(c, http.StatusUnauthorized, nil, "Invalid token claims")
	}

	if claims.UserID == "" {
		logrus.Error("UserID not found in claims")
		return webResponse.ResponseJson(c, http.StatusUnauthorized, nil, "UserID not found in token")
	}

	requestBody := models.UserProfileRequest{
		ID: claims.UserID,
	}

	logrus.Infof("Processing GetProfile | UserID: %s", claims.UserID)

	// Generate Correlation ID
	correlationID := utils.GenerateCorrelationID()

	logrus.Infof("Sending GetProfile event | Correlation ID: %s | UserID: %s", correlationID, claims.UserID)

	err = h.SendMessage.SendingToMessage("GetProfile", correlationID, requestBody)
	if err != nil {
		logrus.Errorf("Failed to send GetProfile message: %v", err)
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to send GetProfile request")
	}

	logrus.Infof("Waiting for GetProfileSuccess/GetProfileFailed response | Timeout: %v", h.Config.RequestTimeout)

	return h.ResponseHandler.HandleEventResponse(
		c,
		false,
		http.StatusOK,
		h.Config.RequestTimeout,
		"Get Profile successfully",
		"GetProfileSuccess",
		"GetProfileFailed",
	)
}
