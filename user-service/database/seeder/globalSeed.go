package seeder

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedAll(db *mongo.Database) {
	SeedUsers(db)
	logrus.Println("Seed all success")
}
