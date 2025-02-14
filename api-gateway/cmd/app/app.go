package main

import (
	"api-gateway/config"
	"api-gateway/handler"
	"api-gateway/routes"
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"os"
	"os/signal"
	"time"
)

type App struct {
	Server  *echo.Echo
	Handler *Handler
}

// Handler Struct for save instance of handler
type Handler struct {
	UserHandler *handler.UserHandler
}

// Initialize prepare environment and setup app
func (app *App) Initialize() {
	app.LoadEnv()
	app.Server = echo.New()

	cfg := config.LoadRateLimitConfig()
	app.Server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(time.Second),
			Burst:     cfg.RateLimit,
			ExpiresIn: 1 * time.Minute,
		},
	)))

	// logging
	config.SetupLogger()
	// recover from panic
	app.Server.Use(middleware.Recover())
	// CORS
	app.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	app.Server.Use(middleware.Gzip())

	routes.UserRoutes(app.Server, cfg)
}

// Run function to run the app
func (app *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// graceful shutdown
	go func() {
		if err := app.Server.Start(":" + port); err != nil {
			logrus.Info("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server shutdown")
}

func (app *App) LoadEnv() {
	appEnv := os.Getenv("APP_ENV")
	envFile := ".env"

	if appEnv == "development" {
		envFile = ".env.development"
	}

	if err := godotenv.Load(envFile); err != nil {
		logrus.Fatal("Error loading .env file ", err)
	}
}
