package repository

import (
	"time"
	"user-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUsersRepository(log *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		Log: log,
	}
}

func (r *UsersRepository) IsExist(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		r.Log.Errorf("Error checking if user exists: %v", err)
		return false, err
	}
	return count > 0, nil
}

func (r *UsersRepository) FindByEmailToken(tx *gorm.DB, token string) (*entity.User, error) {
	var user entity.User
	if err := tx.Where("email_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepository) UpdateEmailVerifiedAt(tx *gorm.DB, userID int) error {
	return tx.Model(&entity.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verified_at": time.Now(),
			"updated_at":        time.Now(),
		}).Error
}

func (r *UsersRepository) FindByEmailVerified(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ? AND email_verified_at IS NOT NULL", email).First(user).Error
}
