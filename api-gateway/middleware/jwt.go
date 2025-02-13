package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

// JWTMiddleware function to check JWT token
func JWTMiddleware(jwtKey []byte) echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey:    jwtKey,
		SigningMethod: "HS256",
		ContextKey:    "user",
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c echo.Context, err error) error {
			// Default status and message
			var status int = http.StatusUnauthorized
			var message string = "Unauthorized"

			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				message = "Token has expired"
				logrus.Warn("JWT Middleware Error: Token expired")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				message = "Invalid token signature"
				logrus.Warn("JWT Middleware Error: Invalid token signature")
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				message = "Token is not yet valid"
				logrus.Warn("JWT Middleware Error: Token is not yet valid")
			case errors.Is(err, jwt.ErrTokenMalformed):
				message = "Malformed token"
				logrus.Warn("JWT Middleware Error: Malformed token")
			default:
				message = "Invalid or missing token"
				logrus.Warn("JWT Middleware Error: Unauthorized access")
			}

			return c.JSON(status, map[string]string{"error": message})
		},
	}
	return echojwt.WithConfig(config)
}
