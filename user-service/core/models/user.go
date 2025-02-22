package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Enum Role USER
const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleUser       = "USER"
)

// User Model struct
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Address   string             `json:"address,omitempty" bson:"address,omitempty"`
	Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Age       int                `json:"age,omitempty" bson:"age,omitempty"`
	GoogleID  string             `json:"google_id,omitempty" bson:"google_id,omitempty"`
	Avatar    string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
