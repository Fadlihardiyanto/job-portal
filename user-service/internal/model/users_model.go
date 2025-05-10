package model

import "time"

type RegisterUserRequest struct {
	ID                   int    `json:"id,omitempty"`
	Password             string `json:"password" validate:"required,gt=6,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	Email                string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID             string    `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	AccessToken    string    `json:"access_token,omitempty"`
	AccessExpiry   time.Time `json:"access_expiry,omitzero"`
	RefreshTokenID string    `json:"refresh_token_id,omitempty"`
	Role           string    `json:"role,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitzero"`
	UpdatedAt      time.Time `json:"updated_at,omitzero"`
}
