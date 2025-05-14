package converter

import (
	"company-service/internal/entity"
	"company-service/internal/model"
)

func CompanyToResponse(company *entity.Company) *model.ResponseCompany {
	return &model.ResponseCompany{
		ID:               company.ID,
		Name:             company.Name,
		City:             company.City,
		OrganizationSize: company.OrganizationSize,
		Logo:             company.Logo,
		UserAccess:       company.UserAccess,
		CreatedAt:        company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func CompanyToEvent(company *entity.Company) *model.CompanyEvent {
	return &model.CompanyEvent{
		ID:               company.ID,
		Name:             company.Name,
		City:             company.City,
		OrganizationSize: company.OrganizationSize,
		Logo:             company.Logo,
		UserAccess:       company.UserAccess,
		CreatedAt:        company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
