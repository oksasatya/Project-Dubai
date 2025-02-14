package seeder

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"user-service/internal/models"
	"user-service/utils"
)

func SeedUsers(db *mongo.Database) {
	userCollection := db.Collection("users")

	var users []interface{}

	for i := 1; i <= 15; i++ {
		gofakeit.Seed(0)
		password, _ := utils.HashPassword("test12345")
		role := gofakeit.RandomString([]string{models.AdminRole, models.UserRole})
		user := models.User{
			Username: gofakeit.Name(),
			Email:    gofakeit.Email(),
			Address:  gofakeit.Address().Address,
			Age:      gofakeit.Number(18, 60),
			Phone:    gofakeit.Phone(),
			Password: password,
			Role:     role,
		}
		users = append(users, user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.InsertMany(ctx, users, options.InsertMany().SetOrdered(false))
	if err != nil {
		logrus.Fatalf("Seed users failed: %v", err)
		return
	}

	logrus.Println("Seed users success")
}
