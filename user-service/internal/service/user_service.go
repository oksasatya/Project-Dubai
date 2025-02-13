package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/pkg"
)

type UserService interface {
	RegisterUser(ctx context.Context, req models.RegisterInput) (*models.User, error)
	LoginUser(ctx context.Context, req *models.LoginInput) (*models.User, error)
}

type userService struct {
	userRepo  repository.UserRepo
	jwtSecret string
}

func (u *userService) RegisterUser(ctx context.Context, req models.RegisterInput) (*models.User, error) {
	hashPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	objectId := primitive.NewObjectID()
	user := models.User{
		ID:       objectId,
		Email:    req.Email,
		Password: hashPassword,
		Address:  req.Address,
		Username: req.Username,
		Age:      req.Age,
		Phone:    req.Phone,
	}

	_, err = u.userRepo.SaveUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userService) LoginUser(ctx context.Context, req *models.LoginInput) (*models.User, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !pkg.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := pkg.GenerateToken(user.ID.Hex(), u.jwtSecret)
	if err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

func NewUserService(userRepo repository.UserRepo, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
