package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserActivityLog struct is used for user activity
type UserActivityLog struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            primitive.ObjectID `bson:"user_id"`
	CourseID          primitive.ObjectID `bson:"course_id,omitempty"`
	ActivityType      string             `bson:"activity_type"`
	ActivityTimestamp primitive.DateTime `bson:"activity_timestamp"`
}
