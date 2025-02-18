package routes

import (
	"api-gateway/config"
	"api-gateway/handler"
	"api-gateway/webResponse"
	"github.com/labstack/echo/v4"
	"messaging"
)

// UserRoutes register user routes
func UserRoutes(e *echo.Echo, cfg *config.RateLimitConfig, rmq *messaging.RabbitMQConnection, res *webResponse.ResponseHandler) {
	userHandler := handler.NewUserHandler(cfg, rmq, res)
	r := e.Group("/api/users")
	//r.POST("/login", userHandler.Login)
	r.POST("/register", userHandler.Register)

	// oauthGroup
	//r.GET("/oauth/google", userHandler.GoogleLogin)
	//r.GET("/oauth/google/callback", userHandler.GoogleCallback)
}
