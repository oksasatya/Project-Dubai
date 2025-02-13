package controller

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) HandleRegistration(message []byte, responseChannel chan string) {
	var req models.RegisterInput
	err := json.Unmarshal(message, &req)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)

		response := models.UserRegistrationResponse{
			Meta: models.MetaDataResponse{
				Message: "Invalid request format",
				Code:    http.StatusBadRequest,
				Status:  "fail",
			},
		}

		responseMsg, _ := json.Marshal(response)
		responseChannel <- string(responseMsg)
		return
	}

	// Call User Service
	user, err := c.userService.RegisterUser(context.Background(), req)
	if err != nil {
		logrus.Errorf("Failed to register user: %v", err)

		response := models.UserRegistrationResponse{
			Meta: models.MetaDataResponse{
				Message: "User registration failed",
				Code:    http.StatusBadRequest,
				Status:  "fail",
			},
		}

		responseMsg, _ := json.Marshal(response)
		responseChannel <- string(responseMsg)
		return
	}

	response := models.UserRegistrationResponse{
		Meta: models.MetaDataResponse{
			Message: "User registered successfully",
			Code:    http.StatusOK,
			Status:  "success",
		},
		Data: models.UserDataResponse{
			Username: user.Username,
			Email:    user.Email,
			Address:  user.Address,
			Age:      user.Age,
			Phone:    user.Phone,
		},
	}

	responseMsg, _ := json.Marshal(response)
	responseChannel <- string(responseMsg)

	logrus.Infof("User registered successfully: %s", user.Username)
}

func (c *UserController) LoginUser(ctx echo.Context) error {

}
