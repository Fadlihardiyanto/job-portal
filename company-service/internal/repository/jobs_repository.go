package repository

import (
	"company-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type JobsRepository struct {
	Repository[entity.Jobs]
	Log *logrus.Logger
}

func NewJobsRepository(log *logrus.Logger) *JobsRepository {
	return &JobsRepository{
		Log: log,
	}
}

func (r *JobsRepository) GetByCompanyID(db *gorm.DB, companyID string) ([]entity.Jobs, error) {
	var jobs []entity.Jobs
	if err := db.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		r.Log.Errorf("Error finding jobs by company ID: %v", err)
		return nil, err
	}
	return jobs, nil
}
