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

type JobsUsecase struct {
	DB           *gorm.DB
	Log          *logrus.Logger
	Validate     *validator.Validate
	Viper        *viper.Viper
	JobsRepo     *repository.JobsRepository
	JobsProducer *messaging.JobsProducer
}

func NewJobsUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, jobsRepo *repository.JobsRepository, jobsProducer *messaging.JobsProducer) *JobsUsecase {
	return &JobsUsecase{
		DB:           db,
		Log:          log,
		Validate:     validate,
		Viper:        viper,
		JobsRepo:     jobsRepo,
		JobsProducer: jobsProducer,
	}
}

func (j *JobsUsecase) CreateJob(ctx context.Context, request *model.RequestJobs) (*model.ResponseJobs, error) {
	tx := j.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := j.Validate.Struct(request)
	if err != nil {
		j.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	job := &entity.Jobs{
		JobsTitle:     request.JobsTitle,
		CompanyID:     request.CompanyID,
		Location:      request.Location,
		WorkspaceType: request.WorkspaceType,
		MinSalary:     request.MinSalary,
		MaxSalary:     request.MaxSalary,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = j.JobsRepo.Create(tx, job)
	if err != nil {
		j.Log.Errorf("Failed to create job: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create job")
	}

	if err := tx.Commit().Error; err != nil {
		j.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	event := converter.JobsToEvent(job)

	if err := j.JobsProducer.Send(event); err != nil {
		j.Log.Errorf("Failed to send job event to Kafka: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send job event to Kafka")
	}

	return converter.JobsToResponse(job), nil
}

func (j *JobsUsecase) GetAllJobs(ctx context.Context) ([]model.ResponseJobs, error) {
	tx := j.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	jobs, err := j.JobsRepo.FindAll(tx)
	if err != nil {
		j.Log.Errorf("Failed to get all jobs: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get all jobs")
	}

	if err := tx.Commit().Error; err != nil {
		j.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	var response []model.ResponseJobs
	for _, job := range jobs {
		response = append(response, *converter.JobsToResponse(&job))
	}

	return response, nil
}

func (j *JobsUsecase) GetJobsByCompanyID(ctx context.Context, request *model.RequestFindJobsByCompanyID) ([]model.ResponseJobs, error) {
	tx := j.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := j.Validate.Struct(request)
	if err != nil {
		j.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	jobs, err := j.JobsRepo.GetByCompanyID(tx, request.CompanyID)
	if err != nil && err != gorm.ErrRecordNotFound {
		j.Log.Errorf("Failed to get jobs by company ID: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get jobs by company ID")
	}

	if err == gorm.ErrRecordNotFound {
		j.Log.Warnf("No jobs found for company ID: %s", request.CompanyID)
		return nil, fiber.NewError(fiber.StatusNotFound, "No jobs found for company ID")
	}

	if err := tx.Commit().Error; err != nil {
		j.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	var response []model.ResponseJobs
	for _, job := range jobs {
		response = append(response, *converter.JobsToResponse(&job))
	}

	return response, nil
}

func (j *JobsUsecase) FindByID(ctx context.Context, request *model.RequestFindJobsByID) (*model.ResponseJobs, error) {
	tx := j.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := j.Validate.Struct(request)
	if err != nil {
		j.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	job, err := j.JobsRepo.FindByID(tx, request.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		j.Log.Errorf("Failed to get job by ID: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get job by ID")
	}

	if err == gorm.ErrRecordNotFound {
		j.Log.Warnf("Job not found with ID: %s", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Job not found")
	}

	if err := tx.Commit().Error; err != nil {
		j.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return converter.JobsToResponse(job), nil
}

func (j *JobsUsecase) UpdateJob(ctx context.Context, request *model.RequestUpdateJobs) (*model.ResponseJobs, error) {
	tx := j.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := j.Validate.Struct(request)
	if err != nil {
		j.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	job, err := j.JobsRepo.FindByID(tx, request.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		j.Log.Errorf("Failed to get job by ID: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get job by ID")
	}

	if err == gorm.ErrRecordNotFound {
		j.Log.Warnf("Job not found with ID: %s", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Job not found")
	}

	job.JobsTitle = request.JobsTitle
	job.CompanyID = request.CompanyID
	job.Location = request.Location
	job.WorkspaceType = request.WorkspaceType
	job.MinSalary = request.MinSalary
	job.MaxSalary = request.MaxSalary
	job.UpdatedAt = time.Now()

	err = j.JobsRepo.Update(tx, job)
	if err != nil {
		j.Log.Errorf("Failed to update job: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update job")
	}

	if err := tx.Commit().Error; err != nil {
		j.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	event := converter.JobsToEvent(job)

	if err := j.JobsProducer.Send(event); err != nil {
		j.Log.Errorf("Failed to send job event to Kafka: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send job event to Kafka")
	}

	return converter.JobsToResponse(job), nil
}
