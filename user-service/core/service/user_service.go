package service

import (
	"context"
	"encoding/json"
	"messaging"
	"time"
	"user-service/api"
	"user-service/core/models"
	"user-service/core/repository"
	"user-service/utils"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	HandleUserRegistered(ctx context.Context, eventData []byte, correlationID string)
	HandleUserLogin(ctx context.Context, eventData []byte, correlationID string)
}

type userService struct {
	userRepo    repository.UserRepo
	rmq         *messaging.RabbitMQConnection
	sendMessage *api.SendingMessage
}

// HandleUserRegistered is a function to handle user registration
func (c *userService) HandleUserRegistered(ctx context.Context, eventData []byte, correlationID string) {
	if c.sendMessage == nil {
		logrus.Fatalf("Failed to initialize SendingMessage")
		return
	}
	// Unmarshal event JSON ke struct `UserRegisteredEvent`
	var req models.UserRegisteredEvent
	if err := json.Unmarshal(eventData, &req); err != nil {
		logrus.Errorf("Invalid event data: %v", err)
		return
	}

	// Cek if user already registered
	existingUser, _ := c.userRepo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		errorResponse := c.sendMessage.SendingToMessage("UserRegisteredFailed", correlationID, "Email already registered")
		if errorResponse != nil {
			logrus.Errorf("Failed to publish UserRegisteredFailed: %v", errorResponse)
		}
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		errorResponse := c.sendMessage.SendingToMessage("UserRegisteredFailed", correlationID, "Failed to hash password")
		if errorResponse != nil {
			logrus.Errorf("Failed to publish UserRegisteredFailed: %v", errorResponse)
		}
		return
	}

	// Simpan user ke database
	newUser := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Address:   req.Address,
		Age:       req.Age,
		Phone:     req.Phone,
		Password:  hashedPassword,
		Role:      models.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = c.userRepo.SaveUser(ctx, &newUser)
	if err != nil {
		errorResponse := c.sendMessage.SendingToMessage("UserRegisteredFailed", correlationID, "Failed to hash password")
		if errorResponse != nil {
			logrus.Errorf("Failed to publish UserRegisteredFailed: %v", errorResponse)
		}
		return
	}

	successResponse := c.sendMessage.SendingToMessage("UserRegisteredSuccess", correlationID, models.UserRegisteredEvent{
		Email:    newUser.Email,
		Username: newUser.Username,
		Address:  newUser.Address,
		Age:      newUser.Age,
		Phone:    newUser.Phone,
		Role:     newUser.Role,
	})
	if successResponse != nil {
		logrus.Errorf("failed to publish User Registerd Success %v", successResponse)
	}
}

// HandleUserLogin is a function to handle user login
func (c *userService) HandleUserLogin(ctx context.Context, eventData []byte, correlationID string) {
	//TODO implement me
	panic("implement me")
}

// NewUserService for handling user service
func NewUserService(userRepo repository.UserRepo, rmq *messaging.RabbitMQConnection) UserService {
	sendMessage := api.NewSendingMessage(rmq)
	return &userService{
		userRepo:    userRepo,
		rmq:         rmq,
		sendMessage: sendMessage,
	}
}
