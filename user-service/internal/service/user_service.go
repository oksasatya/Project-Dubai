package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/utils"
)

type UserService interface {
	RegisterUser(message []byte, responseChannel chan string)
	LoginUser(ctx context.Context, req *models.LoginRequest) (*models.UserLoginResponse, error)
}

type userService struct {
	userRepo  repository.UserRepo
	jwtSecret string
}

func (c *userService) RegisterUser(message []byte, responseChannel chan string) {
	var req models.RegisterRequest
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

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SendErrorResponse(responseChannel, "Error hashing password", http.StatusInternalServerError)
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
		Address:  req.Address,
		Age:      req.Age,
		Phone:    req.Phone,
		Role:     models.UserRole,
	}
	_, err = c.userRepo.SaveUser(context.Background(), &user)
	if err != nil {
		utils.SendErrorResponse(responseChannel, "User registration failed", http.StatusInternalServerError)
		return
	}

	response := models.UserRegistrationResponse{
		Meta: models.MetaDataResponse{
			Message: "User registered successfully",
			Code:    http.StatusCreated,
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

	utils.SendSuccessResponse(responseChannel, response)
}

func (c *userService) LoginUser(ctx context.Context, req *models.LoginRequest) (*models.UserLoginResponse, error) {
	user, err := c.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID.Hex(), c.jwtSecret, user.Role)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &models.UserLoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Address:  user.Address,
		Phone:    user.Phone,
		Age:      user.Age,
		UserToken: struct {
			Token string `json:"token"`
		}{Token: token},
	}, nil
}

func NewUserService(userRepo repository.UserRepo, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
