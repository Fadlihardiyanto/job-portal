package model

import (
	"time"
)

type RegisterUserRequest struct {
	Password             string `json:"password" validate:"required,gt=6,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	Email                string `json:"email" validate:"required,email"`
	Role                 string `json:"role" validate:"required,oneof=admin user"`
}

type UserResponse struct {
	ID             int       `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	AccessToken    string    `json:"access_token,omitempty"`
	AccessExpiry   time.Time `json:"access_expiry,omitzero"`
	RefreshTokenID string    `json:"refresh_token_id,omitempty"`
	Role           string    `json:"role,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitzero"`
	UpdatedAt      time.Time `json:"updated_at,omitzero"`
}

type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Email     string `json:"email" validate:"omitempty,email"`
	About     string `json:"about" validate:"omitempty"`
	Photo     string `json:"photo" validate:"omitempty"`
	Password  string `json:"password" validate:"omitempty,gt=6,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	Role      string `json:"role" validate:"omitempty"`
	IsActive  *bool  `json:"is_active" validate:"omitempty"`
}
