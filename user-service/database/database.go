package database

import (
	"context"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"
	"user-service/database/migrations"
)

var (
	db     *mongo.Database
	client *mongo.Client
	once   sync.Once
)

func InitMongoDB() (*mongo.Database, error) {
	var err error

	once.Do(func() {
		client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		db = client.Database(os.Getenv("MONGO_DB"))

		// migrate
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("Error migrating: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetMongoDB() *mongo.Database {
	return db
}

func CloseMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
