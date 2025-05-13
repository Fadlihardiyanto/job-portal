package repository

import (
	"user-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ResumeRepository struct {
	Repository[entity.Resume]
	Log *logrus.Logger
}

func NewResumeRepository(log *logrus.Logger) *ResumeRepository {
	return &ResumeRepository{
		Log: log,
	}
}

func (r *ResumeRepository) FindResumeByUserID(db *gorm.DB, userID string) ([]entity.Resume, error) {
	var resumes []entity.Resume
	if err := db.Where("user_id = ?", userID).Find(&resumes).Error; err != nil {
		r.Log.Errorf("Error finding resumes by user ID: %v", err)
		return nil, err
	}
	return resumes, nil
}
