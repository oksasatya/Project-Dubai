package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"user-service/internal/models"
)

type UserActivityLogRepo interface {
	CreateUserActivityLog(ctx context.Context, log *models.UserActivityLog) (*mongo.InsertOneResult, error)
	GetUserActivityLogs(ctx context.Context, userId string) ([]*models.UserActivityLog, error)
}

type userActivityLogRepo struct {
	db *mongo.Database
}

func (r *userActivityLogRepo) CreateUserActivityLog(ctx context.Context, log *models.UserActivityLog) (*mongo.InsertOneResult, error) {
	collection := r.db.Collection("UserActivityLog")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, log)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *userActivityLogRepo) GetUserActivityLogs(ctx context.Context, userId string) ([]*models.UserActivityLog, error) {
	collection := r.db.Collection("UserActivityLog")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": objectId}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.UserActivityLog
	for cursor.Next(ctx) {
		var log models.UserActivityLog
		err := cursor.Decode(&log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func NewUserActivityLogRepo(db *mongo.Database) UserActivityLogRepo {
	return &userActivityLogRepo{db: db}
}
