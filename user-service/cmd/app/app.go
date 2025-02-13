package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"messaging"
	"os"
	"os/signal"
	"time"
	"user-service/database"
	"user-service/internal/controller"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg"
)

// App struct for save instance of app
type App struct {
	DB         *mongo.Database
	Server     *echo.Echo
	Controller *Controller
}

type Controller struct {
	UserController *controller.UserController
}

// Initialize prepare environment and setup app
func (app *App) Initialize() {
	app.LoadEnv()

	// setup loger
	pkg.SetupLogger()

	// init db
	db, err := database.InitMongoDB()
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}
	app.DB = db

	// init Controller
	app.Controller = app.InitController()

	// init Server
	app.Server = echo.New()

	// Run consumer
	app.RunConsumer()
}

// RunConsumer function to run consumer
func (app *App) RunConsumer() {
	go messaging.ConsumeMessage("user_registration_queue", app.Controller.UserController.HandleRegistration, messaging.ResponseChannel)
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

// LoadEnv function to load environment variables
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

// InitController function to initialize controller
func (app *App) InitController() *Controller {
	userRepo := repository.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo, os.Getenv("JWT_SECRET"))
	userController := controller.NewUserController(userService)

	return &Controller{
		UserController: userController,
	}
}
