package main

import (
	"api-gateway/config"
	"api-gateway/handler"
	"api-gateway/routes"
	"api-gateway/webResponse"
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"messaging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Server          *echo.Echo
	Handler         *Handler
	DB              *mongo.Database
	RMQ             *messaging.RabbitMQConnection
	ResponseHandler *webResponse.ResponseHandler
}

// Handler Struct for saving instance of handler
type Handler struct {
	UserHandler *handler.UserHandler
}

// Initialize sets up environment and app
func (app *App) Initialize() {
	app.LoadEnv()
	app.Server = echo.New()

	cfg := config.LoadRateLimitConfig()
	app.Server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(cfg.RateLimit),
			Burst:     cfg.RateLimit,
			ExpiresIn: 1 * time.Minute,
		},
	)))

	defer func() {
		if r := recover(); r != nil {
			logrus.Fatalf("Panic occurred during initialization: %v", r)
		}
	}()

	config.SetupLogger()
	app.Server.Use(middleware.Recover())
	app.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))
	app.Server.Use(middleware.Gzip())

	rmq, err := messaging.NewRabbitMQConnection()
	if err != nil {
		logrus.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	app.RMQ = rmq
	logrus.Info("RabbitMQ initialized successfully")

	if app.RMQ == nil {
		logrus.Fatal("Failed to initialize RabbitMQ")
	}
	app.Handler = &Handler{
		UserHandler: handler.NewUserHandler(cfg, app.RMQ, app.ResponseHandler),
	}
	app.ResponseHandler = webResponse.NewResponseHandler(app.RMQ)
	if app.ResponseHandler == nil {
		logrus.Fatal("ResponseHandler is nil after initialization")
	}

	if app.Handler == nil || app.Handler.UserHandler == nil {
		logrus.Fatal("Failed to initialize handler")
	}

	routes.UserRoutes(app.Server, cfg, app.RMQ, app.ResponseHandler)
}

// LoadEnv function to load environment variables
func (app *App) LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file ", err)
	}

	envFile := ".env"

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "development" {
		envFile = ".env.development"
	}

	if err := godotenv.Load(envFile); err != nil {
		logrus.Fatal("Error loading .env file ", err)
	}

	logrus.Printf("Environment running on: %s , .env : %s", appEnv, envFile)
}

// Run starts the server
func (app *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := app.Server.Start(":" + port); err != nil {
			logrus.Fatalf("Server stopped unexpectedly: %v", err)
		}
	}()
	app.handleShutdown()

}

// handleShutdown function to gracefully shutdown server
func (app *App) handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logrus.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Error shutting down server: %v", err)
	}

	if app.RMQ != nil {
		app.RMQ.Close()
	}

	logrus.Info("Server and RabbitMQ connection closed successfully")
}
