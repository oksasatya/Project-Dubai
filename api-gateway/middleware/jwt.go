package middleware

import (
	"api-gateway/config"
	"api-gateway/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var Ctx = context.Background()

// JWTMiddleware function to check JWT token
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logrus.Warn("Missing Authorization header")
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Missing token"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				logrus.Warn("Invalid token format")
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token format"})
			}

			claims := &utils.JWTCustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				logrus.Errorf("Error parsing token: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			if !token.Valid {
				logrus.Warn("Token is not valid")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			// Cek token di Redis
			rdb := config.NewRedisClient()
			ctx := c.Request().Context()

			storedToken, err := rdb.Get(ctx, tokenString).Result()
			if errors.Is(err, redis.Nil) {
				logrus.Warn("Token not found in Redis")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			} else if err != nil {
				logrus.Errorf("Error checking token in Redis: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error accessing Redis"})
			}

			var tokenData map[string]interface{}
			if err := json.Unmarshal([]byte(storedToken), &tokenData); err != nil {
				logrus.Errorf("Error parsing stored token data: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error parsing token data"})
			}

			if claims.UserID != tokenData["userID"].(string) {
				logrus.Warn("Token user ID mismatch")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}
