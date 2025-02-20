package service

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
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
	//HandleUserLogin(ctx context.Context, eventData []byte, correlationID string)
	HandleUserOauth(ctx context.Context, eventData []byte, correlationID string)
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
		errorResponse := c.sendMessage.SendingToMessage("UserRegisteredFailed", correlationID, "Failed to save user")
		if errorResponse != nil {
			logrus.Errorf("Failed to publish UserRegisteredFailed: %v", errorResponse)
		}
		return
	}

	successResponse := c.sendMessage.SendingToMessage("UserRegisteredSuccess", correlationID, map[string]interface{}{
		"email":    newUser.Email,
		"username": newUser.Username,
		"address":  newUser.Address,
		"age":      newUser.Age,
		"phone":    newUser.Phone,
		"role":     newUser.Role,
	})
	if successResponse != nil {
		logrus.Errorf("[RabbitMQ] Failed to publish UserRegisteredSuccess: %v", successResponse)
	}

}

// HandleUserLogin is a function to handle user login
//func (c *userService) HandleUserLogin(ctx context.Context, eventData []byte, correlationID string) {
//	if c.sendMessage == nil {
//		logrus.Fatalf("Failed to initialize SendingMessage")
//		return
//	}
//
//	// Unmarshal event JSON ke struct `UserLoginEvent`
//	var req models.UserLoginEvent
//	if err := json.Unmarshal(eventData, &req); err != nil {
//		logrus.Errorf("Invalid event data: %v", err)
//		return
//	}
//
//	// find user by email
//	user, err := c.userRepo.FindUserByEmail(ctx, req.Email)
//	if err != nil {
//		errorResponse := c.sendMessage.SendingToMessage("UserLoginFailed", correlationID, "Failed to find user")
//		if errorResponse != nil {
//			logrus.Errorf("Failed to publish UserLoginFailed: %v", errorResponse)
//		}
//		return
//	}
//
//	// Cek if password is correct
//	if !utils.CheckPasswordHash(req.Password, user.Password) {
//		errorResponse := c.sendMessage.SendingToMessage("UserLoginFailed", correlationID, "Invalid password")
//		if errorResponse != nil {
//			logrus.Errorf("Failed to publish UserLoginFailed: %v", errorResponse)
//		}
//		return
//	}
//
//	successResponse := c.sendMessage.SendingToMessage("UserLoginSuccess", correlationID, models.UserLoginEvent{
//		Email: user.Email,
//		Role:  user.Role,
//	})
//
//	if successResponse != nil {
//		logrus.Errorf("failed to publish User Login Success %v", successResponse)
//	}
//}

// HandleUserOauth is a function to handle user oauth
func (c *userService) HandleUserOauth(ctx context.Context, eventData []byte, correlationID string) {
	if c.sendMessage == nil {
		logrus.Fatalf("Failed to initialize SendingMessage")
		return
	}

	var req models.UserOAuthEvent

	// Unmarshal event JSON ke struct `UserOauthEvent`
	if err := json.Unmarshal(eventData, &req); err != nil {
		errorResponse := c.sendMessage.SendingToMessage("UserOauthFailed", correlationID, "Failed to unmarshal event data")
		if errorResponse != nil {
			logrus.Errorf("Failed to publish UserOauthFailed: %v", errorResponse)
		}
		return
	}

	// find user by email
	user, err := c.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		// if not found
		if errors.Is(mongo.ErrNoDocuments, err) {
			newUser := models.User{
				GoogleID:  req.GoogleID,
				Email:     req.Email,
				Username:  req.Username,
				Avatar:    req.Avatar,
				Role:      models.RoleUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			_, saveErr := c.userRepo.SaveUser(ctx, &newUser)
			if saveErr != nil {
				logrus.Errorf("Failed to save new user: %v", saveErr)
				errorResponse := c.sendMessage.SendingToMessage("UserOauthFailed", correlationID, "Failed to save user")
				if errorResponse != nil {
					logrus.Errorf("Failed to publish UserOauthFailed: %v", errorResponse)
				}
				return
			}

			//user = savedUser
		} else {
			errorResponse := c.sendMessage.SendingToMessage("UserOauthFailed", correlationID, "Failed to find user")
			if errorResponse != nil {
				logrus.Errorf("Failed to publish UserOauthFailed: %v", errorResponse)
			}
		}

		successResponse := c.sendMessage.SendingToMessage("UserOauthSuccess", correlationID, models.UserOAuthEvent{
			GoogleID: user.GoogleID,
			Email:    user.Email,
			Username: user.Username,
			Avatar:   user.Avatar,
			Role:     user.Role,
		})

		if successResponse != nil {
			logrus.Errorf("failed to publish User Oauth Success %v", successResponse)
		}
	}
}

// NewUserService for handling user service
func NewUserService(userRepo repository.UserRepo, rmq *messaging.RabbitMQConnection, sendMessage *api.SendingMessage) UserService {
	return &userService{
		userRepo:    userRepo,
		rmq:         rmq,
		sendMessage: sendMessage,
	}
}
