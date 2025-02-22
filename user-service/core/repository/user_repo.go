package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"user-service/core/models"
)

type UserRepo interface {
	SaveUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
	SaveToActivityLog(ctx context.Context, activity *models.UserActivityLog) (*mongo.InsertOneResult, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id string) (*models.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*models.User, error)
}

type userRepo struct {
	db *mongo.Database
}

func (r *userRepo) FindByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) SaveUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.Collection("users").FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if errors.Is(mongo.ErrNoDocuments, err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}

	err = r.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) SaveToActivityLog(ctx context.Context, activity *models.UserActivityLog) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.Collection("userActivityLog").InsertOne(ctx, activity)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return &userRepo{db: db}
}
