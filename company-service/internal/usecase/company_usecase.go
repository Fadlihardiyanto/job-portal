package usecase

import (
	"company-service/internal/entity"
	"company-service/internal/gateway/messaging"
	"company-service/internal/model"
	"company-service/internal/model/converter"
	"company-service/internal/repository"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type CompanyUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	Viper             *viper.Viper
	CompanyRepository *repository.CompanyRepository
	CompanyProducer   *messaging.CompanyProducer
}

func NewCompanyUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, companyRepository *repository.CompanyRepository, companyProducer *messaging.CompanyProducer) *CompanyUseCase {
	return &CompanyUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		Viper:             viper,
		CompanyRepository: companyRepository,
		CompanyProducer:   companyProducer,
	}
}

func (c *CompanyUseCase) CreateCompany(ctx context.Context, request *model.RequestCompany) (*model.ResponseCompany, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	company := &entity.Company{
		Name:             request.Name,
		City:             request.City,
		OrganizationSize: request.OrganizationSize,
		Logo:             request.Logo,
		UserAccess:       request.UserAccess,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = c.CompanyRepository.Create(tx, company)
	if err != nil {
		c.Log.Errorf("Failed to create company: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create company")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	event := converter.CompanyToEvent(company)
	if err := c.CompanyProducer.Send(event); err != nil {
		c.Log.Errorf("Failed to send company event to Kafka: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send company event to Kafka")
	}

	return converter.CompanyToResponse(company), nil
}

func (c *CompanyUseCase) GetAllCompany(ctx context.Context) ([]model.ResponseCompany, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	companies, err := c.CompanyRepository.FindAll(tx)
	if err != nil {
		c.Log.Errorf("Failed to get all companies: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get all companies")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	var responseCompanies []model.ResponseCompany
	for _, company := range companies {
		responseCompanies = append(responseCompanies, *converter.CompanyToResponse(&company))
	}

	return responseCompanies, nil
}

func (c *CompanyUseCase) FindByID(ctx context.Context, request *model.RequestFindCompanyByID) (*model.ResponseCompany, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	company, err := c.CompanyRepository.FindByID(tx, request.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Errorf("Failed to get company by ID: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get company by ID")
	}

	if err == gorm.ErrRecordNotFound {
		c.Log.Warnf("Company not found with ID: %s", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Company not found")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return converter.CompanyToResponse(company), nil
}

func (c *CompanyUseCase) UpdateCompany(ctx context.Context, request *model.RequestUpdateCompany) (*model.ResponseCompany, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	company, err := c.CompanyRepository.FindByID(tx, request.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Errorf("Failed to get company by ID: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get company by ID")
	}
	if err == gorm.ErrRecordNotFound {
		c.Log.Warnf("Company not found with ID: %s", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Company not found")
	}

	if request.Name != "" {
		company.Name = request.Name
	}
	if request.City != "" {
		company.City = request.City
	}
	if request.OrganizationSize != "" {
		company.OrganizationSize = request.OrganizationSize
	}
	if request.Logo != "" {
		company.Logo = request.Logo
	}
	if request.UserAccess != "" {
		company.UserAccess = request.UserAccess
	}
	if request.UpdatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, request.UpdatedAt)
		if err != nil {
			c.Log.Warnf("Invalid UpdatedAt format: %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid UpdatedAt format, must be RFC3339")
		}
		company.UpdatedAt = parsedTime
	} else {
		company.UpdatedAt = time.Now()
	}

	err = c.CompanyRepository.Update(tx, company)
	if err != nil {
		c.Log.Errorf("Failed to update company: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update company")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	event := converter.CompanyToEvent(company)
	if err := c.CompanyProducer.Send(event); err != nil {
		c.Log.Errorf("Failed to send company event to Kafka: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send company event to Kafka")
	}

	return converter.CompanyToResponse(company), nil
}

func (c *CompanyUseCase) UpdateCompanyByIDAndUserAccess(ctx context.Context, request *model.RequestUpdateCompanyByIDAndUserAccess) (*model.ResponseCompany, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	company, err := c.CompanyRepository.FindByIDAndUserAccess(tx, request.ID, request.UserAccess)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Errorf("Failed to get company by ID and UserAccess: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get company by ID and UserAccess")
	}
	if err == gorm.ErrRecordNotFound {
		c.Log.Warnf("Company not found with ID: %s and UserAccess: %s", request.ID, request.UserAccess)
		return nil, fiber.NewError(fiber.StatusNotFound, "Company not found")
	}

	if request.Name != "" {
		company.Name = request.Name
	}
	if request.City != "" {
		company.City = request.City
	}
	if request.OrganizationSize != "" {
		company.OrganizationSize = request.OrganizationSize
	}
	if request.Logo != "" {
		company.Logo = request.Logo
	}
	if request.UserAccess != "" {
		company.UserAccess = request.UserAccess
	}
	if request.UpdatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, request.UpdatedAt)
		if err != nil {
			c.Log.Warnf("Invalid UpdatedAt format: %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid UpdatedAt format, must be RFC3339")
		}
		company.UpdatedAt = parsedTime
	} else {
		company.UpdatedAt = time.Now()
	}

	err = c.CompanyRepository.Update(tx, company)
	if err != nil {
		c.Log.Errorf("Failed to update company: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update company")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}
	event := converter.CompanyToEvent(company)
	if err := c.CompanyProducer.Send(event); err != nil {
		c.Log.Errorf("Failed to send company event to Kafka: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send company event to Kafka")
	}

	return converter.CompanyToResponse(company), nil
}
