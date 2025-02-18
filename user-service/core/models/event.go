package models

import "time"

// Event  struct is used for event data
type Event struct {
	EventType     string      `json:"event_type"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"timestamp"`
	Payload       interface{} `json:"payload"`
}

// UserRegisteredEvent struct is used for user registration event
type UserRegisteredEvent struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
}

// UserLoginSuccessEvent User Login Success Event
type UserLoginSuccessEvent struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserOAuthEvent User Register via OAuth
type UserOAuthEvent struct {
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UserOAuthSuccessEvent User Register OAuth Success
type UserOAuthSuccessEvent struct {
	Email string `json:"email"`
}
