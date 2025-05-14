package usecase

import (
	"context"
	"reflect"
	"time"
	common "user-service/internal/common/error"
	commonUtil "user-service/internal/common/util"
	"user-service/internal/gateway/messaging"
	"user-service/internal/model"
	"user-service/internal/model/converter"
	"user-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ResumeUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	Viper            *viper.Viper
	ResumeRepository *repository.ResumeRepository
	ResumeProducer   *messaging.ResumeProducer
}

func NewResumeUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper, resumeRepository *repository.ResumeRepository, resumeProducer *messaging.ResumeProducer) *ResumeUseCase {
	return &ResumeUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		Viper:            viper,
		ResumeRepository: resumeRepository,
		ResumeProducer:   resumeProducer,
	}
}

func (c *ResumeUseCase) CreateResume(ctx context.Context, request *model.RequestResume) (*model.ResponseResume, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	event := &model.ResumeEvent{
		Name:       request.Name,
		Attachment: request.Attachment,
		UserID:     request.UserID,
		Status:     "queued",
		CreatedAt:  time.Now().Format(time.RFC3339),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}

	err = c.ResumeProducer.Send(event)
	if err != nil {
		c.Log.Errorf("Failed to send resume event to Kafka: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to queue resume for processing")
	}

	return converter.EventToResponse(event), nil
}

func (c *ResumeUseCase) GetAllResume(ctx context.Context) ([]model.ResponseResume, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	resume, err := c.ResumeRepository.FindAll(tx)
	if err != nil {
		c.Log.Errorf("Failed to get all resume: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get all resume")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	response := make([]model.ResponseResume, len(resume))
	for i, r := range resume {
		response[i] = *converter.ResumeToResponse(&r)
	}

	return response, nil
}

func (c *ResumeUseCase) FindByID(ctx context.Context, request *model.RequestFindResumeByID) (*model.ResponseResume, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	resume, err := c.ResumeRepository.FindByID(tx, request.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Log.Errorf("Failed to get resume by ID: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get resume by ID")
	}

	if err == gorm.ErrRecordNotFound {
		c.Log.Warnf("Resume not found with ID: %s", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Resume not found")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	response := converter.ResumeToResponse(resume)

	return response, nil
}

func (c *ResumeUseCase) GetByUserID(ctx context.Context, request *model.RequestFindResumeByUser) ([]model.ResponseResume, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		validationErrors := common.FormatValidationError(err, reflect.TypeOf(*request))
		return nil, fiber.NewError(fiber.StatusBadRequest, commonUtil.MapToJSON(validationErrors))
	}

	resume, err := c.ResumeRepository.FindResumeByUserID(tx, request.UserID)
	if err != nil {
		c.Log.Errorf("Failed to get resume by user ID: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get resume by user ID")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	response := make([]model.ResponseResume, len(resume))
	for i, r := range resume {
		response[i] = *converter.ResumeToResponse(&r)
	}

	return response, nil
}
