package converter

import (
	"company-service/internal/entity"
	"company-service/internal/model"
	"time"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToLoginResponse(user *entity.User, accessToken string, accessExpiry time.Time) *model.UserResponse {
	return &model.UserResponse{
		AccessToken:  accessToken,
		AccessExpiry: accessExpiry,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func UserToEvent(user *entity.User) *model.UserEvent {
	return &model.UserEvent{
		ID:        user.ID,
		Email:     user.Email,
		Token:     user.Token, // Token diambil dari parameter, bukan dari entity
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
