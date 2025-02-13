package routes

import (
	"api-gateway/config"
	"api-gateway/models"
	"api-gateway/utils"
	"api-gateway/webResponse"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"messaging"
	"net/http"
	"time"
)

// RegisterUserRoutes register user routes
func RegisterUserRoutes(e *echo.Echo) {
	r := e.Group("/api/users")

	r.POST("/login", func(c echo.Context) error {
		return config.ForwardProxy(c, "http://user-service:8081/api/users/login")
	})

	r.POST("/register", func(c echo.Context) error {
		var requestBody models.RegisterRequest
		if err := c.Bind(&requestBody); err != nil {
			return webResponse.ResponseJson(c, http.StatusBadRequest, nil, "Invalid request format")
		}

		if err := requestBody.Validate(); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				formattedError := utils.FormatValidationError(&requestBody, validationErrors)
				return webResponse.ResponseJson(c, http.StatusBadRequest, formattedError, "Validation failed")
			}
			return webResponse.ResponseJson(c, http.StatusBadRequest, nil, err.Error())
		}

		message, err := json.Marshal(requestBody)
		if err != nil {
			return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to convert request body to JSON")
		}

		err = messaging.PublishMessage("user_registration_queue", string(message))
		if err != nil {
			return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to publish message")
		}

		select {
		case response := <-messaging.ResponseChannel:
			var jsonResponse map[string]interface{}
			if err := json.Unmarshal([]byte(response), &jsonResponse); err != nil {
				logrus.Errorf("Failed to unmarshal response: %v", err)
				return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to process response")
			}
			return c.JSON(http.StatusOK, jsonResponse)

		case <-time.After(10 * time.Second):
			return webResponse.ResponseJson(c, http.StatusGatewayTimeout, nil, "Request timed out waiting for response")
		}
	})

}
