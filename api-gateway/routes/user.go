package routes

import (
	"api-gateway/config"
	"api-gateway/handler"
	"github.com/labstack/echo/v4"
)

// UserRoutes register user routes
func UserRoutes(e *echo.Echo, cfg *config.RateLimitConfig) {
	userHandler := handler.NewUserHandler(*cfg)
	r := e.Group("/api/users")
	// Auth routes
	r.POST("/login", userHandler.Login)
	r.POST("/register", userHandler.Register)
}
