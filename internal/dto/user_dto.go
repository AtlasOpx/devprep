package dto

import (
	"github.com/google/uuid"
	"time"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	Username  string `json:"username,omitempty" validate:"omitempty,min=3,max=100"`
}

type UserProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
}

type UsersListResponse struct {
	Users []UserProfileResponse `json:"users"`
	Total int                   `json:"total"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
