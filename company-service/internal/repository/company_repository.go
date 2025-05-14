package repository

import (
	"company-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	Repository[entity.Company]
	Log *logrus.Logger
}

func NewCompanyRepository(log *logrus.Logger) *CompanyRepository {
	return &CompanyRepository{
		Log: log,
	}
}

func (c *CompanyRepository) FindByIDAndUserAccess(db *gorm.DB, id string, userAccess string) (*entity.Company, error) {
	var company entity.Company
	err := db.Where("id = ? AND user_access = ?", id, userAccess).First(&company).Error
	if err != nil {
		c.Log.Errorf("Error finding company by ID and user access: %v", err)
		return nil, err
	}
	return &company, nil
}
