package usecase

import (
	"company-service/internal/gateway/messaging"
	"company-service/internal/model"
	"company-service/internal/model/converter"
	"company-service/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserJobsUsecase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	Viper            *viper.Viper
	UserJobsRepo     *repository.UserJobsRepository
	JobsRepo         *repository.JobsRepository
	UserJobsProducer *messaging.UserJobsProducer
}

func NewUserJobsUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, userJobsRepo *repository.UserJobsRepository, jobsRepo *repository.JobsRepository, userJobsProducer *messaging.UserJobsProducer) *UserJobsUsecase {
	return &UserJobsUsecase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		Viper:            viper,
		UserJobsRepo:     userJobsRepo,
		JobsRepo:         jobsRepo,
		UserJobsProducer: userJobsProducer,
	}
}

func (u *UserJobsUsecase) CreateUserJobs(ctx context.Context, request *model.UserJobsRequest) (*model.UserJobsResponse, error) {

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Check if the job exists
	job, err := u.JobsRepo.FindByID(u.DB, request.JobID)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Log.Warnf("Job not found: %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Job not found")
	}

	if job == nil {
		u.Log.Warnf("Job not found: %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Job not found")
	}

	event := &model.UserJobsEvent{
		UserID:   request.UserID,
		JobID:    request.JobID,
		ResumeID: request.ResumeID,
		Status:   "queued",
	}

	err = u.UserJobsProducer.Send(event)
	if err != nil {
		u.Log.Errorf("Failed to send user jobs event to Kafka: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to send user jobs event")
	}

	return converter.EventToResponseUserJobs(event), nil
}
