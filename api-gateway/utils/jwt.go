package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateToken generates a JWT token
func GenerateToken(userId string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
