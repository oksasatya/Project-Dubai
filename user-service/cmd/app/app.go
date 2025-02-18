package main

import (
	"api-gateway/config"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"messaging"
	"os"
	"os/signal"
	"sync"
	"time"
	"user-service/core/models"
	"user-service/core/repository"
	"user-service/core/service"
	"user-service/database"
)

// App struct for save instance of app
type App struct {
	DB      *mongo.Database
	Server  *echo.Echo
	Service *Service
	RMQ     *messaging.RabbitMQConnection
}

type Service struct {
	UserService service.UserService
}

// Initialize prepare environment and setup app
func (app *App) Initialize() {
	app.LoadEnv()

	// Init Server
	app.Server = echo.New()

	// Init Database
	db, err := database.InitMongoDB()
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}
	app.DB = db

	config.SetupLogger()

	// Init RabbitMQ
	rmq, err := messaging.NewRabbitMQConnection()
	if err != nil {
		logrus.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	app.RMQ = rmq

	// Init Service
	app.Service = &Service{
		UserService: service.NewUserService(repository.NewUserRepo(db), rmq),
	}
}

// RunConsumer function to run consumer
func (app *App) RunConsumer(wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		messaging.ConsumeEvent(app.RMQ, "UserRegistered", func(event models.Event) {
			ctx := context.Background()

			var req models.UserRegisteredEvent
			payloadBytes, _ := json.Marshal(event.Payload)
			if err := json.Unmarshal(payloadBytes, &req); err != nil {
				logrus.Errorf("Failed to parse event payload: %v", err)
				return
			}

			logrus.Infof("[user-service] Processing UserRegistered | Email: %s", req.Email)
			app.Service.UserService.HandleUserRegistered(ctx, payloadBytes, event.CorrelationID)
		})
	}()

	// **Don't let the main goroutine finish so the consumer will keep running**
	select {}
}

// Run function to run the app
func (app *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// Run Consumer
	go app.RunConsumer(&wg)

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

	if app.RMQ != nil {
		app.RMQ.Close()
	}

	logrus.Info("Server shutdown")
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
