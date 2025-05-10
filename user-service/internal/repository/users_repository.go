package repository

import (
	"user-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersRepository struct {
	Repository[entity.Users]
	Log *logrus.Logger
}

func NewUsersRepository(log *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		Log: log,
	}
}

func (r *UsersRepository) IsExist(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Model(&entity.Users{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		r.Log.Errorf("Error checking if user exists: %v", err)
		return false, err
	}
	return count > 0, nil
}
