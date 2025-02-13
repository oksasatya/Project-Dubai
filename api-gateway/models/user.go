package models

import (
	"github.com/go-playground/validator/v10"
)

// RegisterRequest struct for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Address  string `json:"address" validate:"required"`
	Age      int    `json:"age" validate:"required,gte=18"`
	Phone    string `json:"phone" validate:"required,e164"`
}

// Validate function to validate RegisterRequest using go-playground/validator
func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
