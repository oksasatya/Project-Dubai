package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Variable to store user role
const (
	SuperAdminRole = "SUPER ADMIN"
	AdminRole      = "ADMIN"
	UserRole       = "USER"
)

// User struct is used for user model
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Address  string             `json:"address" bson:"address"`
	Phone    string             `json:"phone" bson:"phone"`
	Age      int                `json:"age" bson:"age"`
	Token    string             `json:"token" bson:"token"`
	Role     string             `json:"role" bson:"role"`
}

// RegisterRequest struct is used for user registration input
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Address  string `json:"address" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Age      int    `json:"age" validate:"required,gt=0"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Validate function to validate LoginRequest using go-playground/validator
func (l *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

// UserLoginResponse struct is used for user response
type UserLoginResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Age       int    `json:"age"`
	UserToken struct {
		Token string `json:"token"`
	}
}

type MetaDataResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type UserDataResponse struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Address  string `json:"address,omitempty"`
	Age      int    `json:"age,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type UserRegistrationResponse struct {
	Meta MetaDataResponse `json:"meta"`
	Data UserDataResponse `json:"data"`
}
