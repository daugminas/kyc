package domain

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	// Id        primitive.ObjectID `json:"id,omitempty" bson:"id"`
	Email     string    `json:"email" validate:"required,email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Active    bool      `json:"active" validate:"boolean" bson:"active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	// UserName  string             `json:"username" validate:"required" bson:"username"`
}

type UserResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}
