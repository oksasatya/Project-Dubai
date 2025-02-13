package migrations

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration is a struct to define migration
type Migration struct {
	ID       string
	Migrate  func() error
	Rollback func() error
}

// Migration function for create_users_collection
func createUsersCollectionMigration(database *mongo.Database, indexField string) *Migration {
	return &Migration{
		ID: "20240909202112_create_users_collection",
		Migrate: func() error {
			collection := database.Collection("users")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Dynamic index options for the given field
			indexOptions := options.Index().SetUnique(true)
			indexModel := mongo.IndexModel{
				Keys:    bson.M{indexField: 1}, // Dynamic field for index creation
				Options: indexOptions,
			}

			_, err := collection.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				return err
			}

			logrus.Printf("Migration: %s completed. Index created on field: %s", "create_users_collection", indexField)
			return nil
		},
		Rollback: func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := database.Collection("users").Drop(ctx)
			if err != nil {
				return err
			}

			logrus.Printf("Rollback: %s completed", "create_users_collection")
			return nil
		},
	}
}
